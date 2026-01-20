package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// opening RabbitMQ connection
	amqpConn, err := amqp.Dial(pubsub.RabbitUrl)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	amqpChan, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("could not create a channel: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		amqpConn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix,
		routing.ArmyMovesPrefix+".*",
		pubsub.Transient,
	)
	if err != nil {
		log.Fatalf("could not declare and bind: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	err = pubsub.SubscribeGob(
		amqpConn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.Durable,
		handlerWriteLogs,
	)
	if err != nil {
		log.Fatalf("could not declare and bind: %v", err)
	}

	gamelogic.PrintServerHelp()

loop:
	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}

		switch userInput[0] {
		case "pause":
			log.Println("sending a pause message")
			pubsub.PublishJSON(amqpChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
		case "resume":
			log.Println("sending a resume message")
			pubsub.PublishJSON(amqpChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
		case "quit":
			log.Println("exiting...")
			break loop
		default:
			log.Println("command is not supported")
		}
	}

	signalChan := make(chan os.Signal, 1)
	// this wait for ctrl+c
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
