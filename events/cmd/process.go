package cmd

import (
	"encoding/json"
	"log"

	"events/pkg/publisher/kafka"
	"events/proto/protos"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var processCmd = &cobra.Command{
	Use: "process",
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

			return queue.Consume(h.process, inputTopic)
		})

		if err = eg.Wait(); err != nil {
			return err
		}

		return nil
	},
}

type Account struct {
	ID       string `json:"account_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *handler) process(msg []byte) error {
	acc := &Account{}
	if err := json.Unmarshal(msg, acc); err != nil {
		return err
	}

	log.Printf("consumed application for accountID: %s", acc.ID)

	event := &protos.AccountCreatedEvent{
		AccountId: acc.ID,
	}

	payload, err := types.MarshalAny(event)
	if err != nil {
		return err
	}

	env := &protos.Envelope{
		Payload: payload,
	}

	data, err := env.Marshal()
	if err != nil {
		return err
	}

	err = h.queue.Publish(data, outputTopic)
	if err != nil {
		return errors.Wrap(err, "publishing")
	}

	return nil
}
