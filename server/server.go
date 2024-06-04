package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"rabbitmq-wrapper/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Missing config file argument")
	}

	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	sc, err := config.ParseConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sc)

	conn, err := amqp.Dial(sc.Url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	http.HandleFunc("POST /publish", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		bytebody, err := io.ReadAll(r.Body)
		fmt.Println(string(bytebody))
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

		if err := json.Unmarshal(bytebody, &messageRequest); err != nil {
			io.WriteString(w, "error")
			failOnError(err, "Failed to decode json")
		}

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
		}

		messageResponse.Body = messageRequest.Body
		messageResponse.UserId = messageRequest.UserId

		messagePublish, err := json.Marshal(messageResponse)
		if err != nil {
			io.WriteString(w, "error")
			failOnError(err, "Failed to marshal messageResponse")
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
		}

		log.Printf(" [x] Sent %s\n", messagePublish)
	})

	err = http.ListenAndServe(sc.Port, nil)
	failOnError(err, "Failed to serve http")
}
