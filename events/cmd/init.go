package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"events/pkg/publisher/kafka"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

var (
	addresses   []string
	inputTopic  string
	outputTopic string
)

var initCmd = &cobra.Command{
	Use: "init",
	Long: `Init is the command that stores initial account inforamtion contained in a http request
	and publishes json object to the initial queue where actual event processing begins. Why does it publish whole request and not AccountCreatedEvent? Well it could all depending on requirements and how many consumers/proicesses there are going to be with requirement for initial request data and not just accountID`,
	PreRunE: func(*cobra.Command, []string) error {
		var address string
		if address = viper.GetString(kafkaAddressFlag); address == "" {
			return errors.New("address must not be empty")
		}

		addresses = append(addresses, address)

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

		db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
		if err != nil {
			return err
		}

		h := &handler{
			queue: queue,
			db:    db,
		}

		http.HandleFunc("/create", h.handle)

		port := 8080

		log.Println("Serving on port", port)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

		return nil
	},
}

type handler struct {
	queue *kafka.Queue
	db    *badger.DB
}

func (h *handler) handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err) // we all know about panics so no negative feedback please
	}

	txn := h.db.NewTransaction(true)
	defer txn.Discard()

	accountID := uuid.New().String()

	err = txn.Set([]byte(accountID), body)
	if err != nil {
		panic(err)
	}

	err = txn.Commit()
	if err != nil {
		panic(err)
	}

	payload, err := sjson.Set(string(body), "account_id", accountID)
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("stored in db for accountID: (%s)", accountID))

	err = h.queue.Publish([]byte(payload), outputTopic)
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("published for accountID: (%s)", accountID))
}
