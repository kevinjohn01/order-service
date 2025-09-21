package queue

import (
    "encoding/json"
    "log"
    "github.com/streadway/amqp"
)

var Conn *amqp.Connection
var Ch *amqp.Channel

func InitRabbitMQ() {
    var err error
    Conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
    if err != nil {
        log.Fatal("Failed to connect to RabbitMQ:", err)
    }

    Ch, err = Conn.Channel()
    if err != nil {
        log.Fatal("Failed to open a channel:", err)
    }

    _, err = Ch.QueueDeclare("order.created", true, false, false, false, nil)
    if err != nil {
        log.Fatal("Failed to declare queue:", err)
    }
}

func PublishOrder(order interface{}) {
    body, _ := json.Marshal(order)
    err := Ch.Publish(
        "", "order.created", false, false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    if err != nil {
        log.Println("Failed to publish:", err)
    }
}
