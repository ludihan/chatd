package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
    config, err := os.ReadFile("/home/lucca/projects/chatd/server_config.toml")
    if err != nil {
        fmt.Println(err)
        return
    }

    sc, err := ParseConfigFile(config)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(sc)
    return


	URL := os.Getenv("AMQP_URL")
	if URL == "" {
		fmt.Println("Could not find value for enviroment variable \"AMQP_URL\"")
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	conn, err := amqp.Dial(URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GOROUTINES", runtime.NumGoroutine())
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		defer r.Body.Close()

		messageRequest := struct {
			Exchange string `json:"exchange"`
			Body     string `json:"body"`
			UserId   string `json:"userId"`
		}{}
		messageResponse := struct {
			Body   string `json:"body"`
			UserId string `json:"userId"`
		}{}

		err := d.Decode(&messageRequest)
		if err != nil {
			io.WriteString(w, "error")
			failOnError(err, "Failed to decode json")
			return
		}

		fmt.Println("exchange:", messageRequest.Exchange)

		err = ch.ExchangeDeclare(
			messageRequest.Exchange, // name
			"fanout",                // type
			false,                   // durable
			true,                    // auto-deleted
			false,                   // internal
			false,                   // no-wait
			nil,                     // arguments
		)
		if err != nil {
			io.WriteString(w, "error")
			failOnError(err, "Failed to declare an exchange")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		messageResponse.Body = messageRequest.Body
		messageResponse.UserId = messageRequest.UserId

		messagePublish, err := json.Marshal(messageResponse)
		if err != nil {
			io.WriteString(w, "error")
			failOnError(err, "Failed to marshal messageResponse")
			return
		}

		err = ch.PublishWithContext(ctx,
			messageRequest.Exchange, // exchange
			"",                      // routing key
			false,                   // mandatory
			false,                   // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        messagePublish,
			})
		if err != nil {
			io.WriteString(w, "error")
			failOnError(err, "Failed to publish a message")
			return
		}

		log.Printf(" [x] Sent %s\n", messagePublish)
	})

    err = http.ListenAndServe(":8080", nil)
	failOnError(err, "Failed to serve http")
}
