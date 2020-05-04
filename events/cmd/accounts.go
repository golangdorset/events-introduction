package cmd

import (
	"log"

	"events/pkg/publisher/kafka"
	"events/proto/protos"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var accountsCmd = &cobra.Command{
	Use: "accounts",
	PreRunE: func(*cobra.Command, []string) error {
		var address string
		if address = viper.GetString(kafkaAddressFlag); address == "" {
			return errors.New("address must not be empty")
		}

		addresses = append(addresses, address)

		if inputTopic = viper.GetString(inputTopicFlag); inputTopic == "" {
			return errors.New("input topic must not be empty")
		}

		if outputTopic = viper.GetString(outputTopicFlag); outputTopic == "" {
			return errors.New("output topic must not be empty")
		}

		return nil
	},
	RunE: func(*cobra.Command, []string) error {
		queue, err := kafka.New(addresses)
		if err != nil {
			return err
		}
		log.Println("kafka connected")

		var eg errgroup.Group
		eg.Go(func() error {
			h := &handler{
				queue: queue,
			}

			return queue.Consume(h.processAccount, inputTopic)
		})

		if err = eg.Wait(); err != nil {
			return err
		}

		return nil
	},
}

func (h *handler) processAccount(msg []byte) error {
	log.Println("consumed event")
	env := &protos.Envelope{}
	if err := proto.Unmarshal(msg, env); err != nil {
		return errors.Wrap(err, "unmarshal envelope")
	}

	switch env.Payload.TypeUrl {
	case "type.googleapis.com/AccountCreatedEvent":
		eventPayload := &protos.AccountCreatedEvent{}
		err := types.UnmarshalAny(env.Payload, eventPayload)
		if err != nil {
			return errors.Wrap(err, "unmarshal payload")
		}

		log.Printf("consumed AccountCreatedEvent for accountID: %s", eventPayload.AccountId)

		changed := &protos.AccountValidatedEvent{
			AccountId: eventPayload.AccountId,
		}

		payload, err := types.MarshalAny(changed)
		if err != nil {
			return err
		}

		enve := &protos.Envelope{
			Payload: payload,
		}

		data, err := enve.Marshal()
		if err != nil {
			return err
		}

		err = h.queue.Publish(data, outputTopic)
		if err != nil {
			return errors.Wrap(err, "publishing")
		}

		return nil
	case "type.googleapis.com/AccountValidatedEvent":
		eventPayload := &protos.AccountValidatedEvent{}
		err := types.UnmarshalAny(env.Payload, eventPayload)
		if err != nil {
			return errors.Wrap(err, "unmarshal payload")
		}

		log.Printf("consumed AccountValidatedEvent for accountID: %s", eventPayload.AccountId)

		// preform some action based on requirements

		return nil
	case "type.googleapis.com/PasswordChangedEvent":
		eventPayload := &protos.PasswordChangedEvent{}
		err := types.UnmarshalAny(env.Payload, eventPayload)
		if err != nil {
			return errors.Wrap(err, "unmarshal payload")
		}

		log.Printf("consumed PasswordChangedEvent for accountID: %s", eventPayload.AccountId)

		// preform action like update db for example

		return nil

	default:
		// other events we are not interested in so return nil to allow further processing and not block the queue
		return nil
	}
}
