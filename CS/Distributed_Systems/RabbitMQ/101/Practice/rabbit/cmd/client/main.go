package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	// opening RabbitMQ connection
	conn, err := amqp.Dial(pubsub.RabbitUrl)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	gs := gamelogic.NewGameState(username)

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create a channel: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.Transient,
		handlerMove(publishCh, gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to army moves: %v", err)
	}
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		pubsub.Durable,
		handlerWar(publishCh, gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to war declarations: %v", err)
	}
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

loop:
	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}

		switch userInput[0] {
		case "spawn":
			// example: "spawn europe infantry"
			gs.CommandSpawn(userInput)
		case "move":
			// example:  "move europe 1"
			// move their units to a new location
			mv, err := gs.CommandMove(userInput)
			if err != nil {
				log.Fatalf("command move failed: %v", err)
			}

			// Publish the move to the army_moves.username routing key
			pubsub.PublishJSON(publishCh, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+gs.GetUsername(), mv)
			log.Println("move was published successfully")
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(userInput) < 2 {
				fmt.Println("usage: spam <n>")
				continue
			}
			n, err := strconv.Atoi(userInput[1])
			if err != nil {
				fmt.Printf("error: %s is not a valid number\n", userInput[1])
				continue
			}
			for i := 0; i < n; i++ {
				msg := gamelogic.GetMaliciousLog()
				err = PublishGameLogs(publishCh, routing.GameLog{
					CurrentTime: time.Now(),
					Username:    gs.GetUsername(),
					Message:     msg,
				}, gs.GetUsername())
				if err != nil {
					fmt.Printf("error publishing malicious log: %s\n", err)
				}
			}
			fmt.Printf("Published %v malicious logs\n", n)
		case "quit":
			gamelogic.PrintQuit()
			break loop
		default:
			log.Println("command is not supported")
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Good Bye from Client...")
}

func PublishGameLogs(amqpChan *amqp.Channel, gameLog routing.GameLog, warInitiator string) error {
	return pubsub.PublishGob(amqpChan, routing.ExchangePerilTopic, routing.GameLogSlug+"."+warInitiator, gameLog)
}
