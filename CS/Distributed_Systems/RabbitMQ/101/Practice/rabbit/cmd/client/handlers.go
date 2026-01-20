package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(amqpChan *amqp.Channel, gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(mv gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(mv)
		switch moveOutcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(amqpChan, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix+"."+gs.GetUsername(), gamelogic.RecognitionOfWar{
				Attacker: mv.Player,
				Defender: gs.GetPlayerSnap(),
			})
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		default:
			return pubsub.NackDiscard
		}
	}
}

func handlerWar(amqpChan *amqp.Channel, gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			err := PublishGameLogs(amqpChan, routing.GameLog{
				CurrentTime: time.Now(),
				Message:     fmt.Sprintf("%s won a war against %s", winner, loser),
				Username:    rw.Attacker.Username,
			}, rw.Attacker.Username)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			err := PublishGameLogs(amqpChan, routing.GameLog{
				CurrentTime: time.Now(),
				Message:     fmt.Sprintf("%s won a war against %s", winner, loser),
				Username:    rw.Attacker.Username,
			}, rw.Attacker.Username)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			err := PublishGameLogs(amqpChan, routing.GameLog{
				CurrentTime: time.Now(),
				Message:     fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser),
				Username:    rw.Attacker.Username,
			}, rw.Attacker.Username)

			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			PublishGameLogs(amqpChan, routing.GameLog{
				CurrentTime: time.Now(),
				Message:     fmt.Sprintf("unexpected war outcome: %s", outcome),
				Username:    rw.Attacker.Username,
			}, rw.Attacker.Username)
			return pubsub.NackDiscard
		}
	}
}
