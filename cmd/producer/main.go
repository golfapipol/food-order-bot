package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("connection failed", err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("open channel failed", err.Error())
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"foodbot", // name
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal("declare queue failed", err.Error())
	}

	engine := gin.Default()

	engine.POST("/webhook", func(context *gin.Context) {
		var request gin.H
		if err := context.ShouldBindJSON(&request); err != nil {
			context.Error(err)
			return
		}
		bytes, err := json.Marshal(request)
		if err != nil {
			context.Error(err)
			return

		}
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        bytes,
			})
		log.Printf(" [x] Sent %s", bytes)
		if err != nil {
			log.Fatal("register failed ", err.Error())
		}
		context.JSON(200, gin.H{})
	})

	engine.Run(":3000")

}
