package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
    if err != nil {
        log.Panicf("%s: %s", msg, err)
    }
}

func main() {
    var exchange string
    if len(os.Args) >= 2 {
        exchange = os.Args[1]
    } else {
        fmt.Println("no args")
        return
    }

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

    err = ch.ExchangeDeclare(
        exchange, // name
        "fanout", // type
        false,    // durable
        true,     // auto-deleted
        false,    // internal
        false,    // no-wait
        nil,      // arguments
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
        q.Name,   // queue name
        "",       // routing key
        exchange, // exchange
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
                Body string `json:"body"`
                UserId string `json:"userId"`
                Exchange string `json:"exchange"`
            }{}

            err = json.Unmarshal(d.Body, &recMessage)
            if err != nil {
                fmt.Println("unmarshall deu errado")
            }
            fmt.Println(recMessage)

            f.Write([]byte(fmt.Sprintf("[%v] %v\n", recMessage.UserId, recMessage.Body)))
        }
    }()

    go func() {
        for {
            fileInputInit(f, ch)
        }
    }()

    log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
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
        Body string `json:"body"`
        UserId string `json:"userId"`
        Exchange string `json:"exchange"`
    }{Exchange: os.Args[1], UserId: os.Args[2], Body: input}

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    mensagem, _ := json.Marshal(recMessage)

    ch.PublishWithContext(ctx,
    os.Args[1], // exchange
    "",                      // routing key
    false,                   // mandatory
    false,                   // immediate
    amqp.Publishing{
        ContentType: "application/json",
        Body:        mensagem,
    })
}
