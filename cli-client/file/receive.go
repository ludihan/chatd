package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"rabbitmq-wrapper/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
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

	err = ch.ExchangeDeclare(
		sc.Exchange, // name
		"fanout",    // type
		false,       // durable
		true,        // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,      // queue name
		"",          // routing key
		sc.Exchange, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	f, err := os.Create("chat.log")

	go func() {
		for d := range msgs {
			recMessage := struct {
				Body     string `json:"body"`
				UserId   string `json:"userId"`
				Exchange string `json:"exchange"`
			}{}

			err = json.Unmarshal(d.Body, &recMessage)
			if err != nil {
				fmt.Println("unmarshall deu errado")
			}

			f.Write([]byte(fmt.Sprintf("[%v] %v\n", recMessage.UserId, recMessage.Body)))
		}
	}()

	go func() {
		for {
			fileInputInit(f, ch)
		}
	}()

	<-forever
}

func fileInputInit(file *os.File, ch *amqp.Channel) {
	var input string
	fmt.Scanln(&input)

	if input == "exit" {
		file.Close()
		os.Exit(0)
	}

	recMessage := struct {
		Body     string `json:"body"`
		UserId   string `json:"userId"`
		Exchange string `json:"exchange"`
	}{Exchange: os.Args[1], UserId: os.Args[2], Body: input}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mensagem, _ := json.Marshal(recMessage)

	ch.PublishWithContext(ctx,
		os.Args[1], // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        mensagem,
		})
}
