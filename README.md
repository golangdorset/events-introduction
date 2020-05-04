# Introduction to events

This is my first presentation about event streamming/sourcing/DDD.

## Steps

1. Start kafka in docker

       docker-compose -f events/resources/kafka/docker-compose.yml up -d

1. Build the binary

       go build -o events

1. Start http server accepting account creation  

       ./events init --kafka-address localhost:9092 --output-topic accounts

1. We are going to assume all account related events will be published in `accounts` topic. One topic can store many types of events as long as they are going to be wrapped in [Envelope](events/proto/envelope.proto) event type where payload will be strictly related event type.

1. Make http request which will store data in database (badger in this case) and publish this data to the `init-requests` topic for further processing or data availability by other domains so they dont have to query the database each time.

       curl -XPOST -H "Content-Type: application/json" localhost:8080/create -d "{\"name\": \"gopher\", \"password\": \"v3ry-s3cr3t\"}"

1. Start a consumer that consumes from `init-requests` topic and publishes `AccountCreatedEvent` to `accounts` topic.

       ./events process --kafka-address localhost:9092 --input-topic init --output-topic accounts

1. Start a consumer that consumes from `accounts` topic and publishes many different events to the same topic.

       ./events process --kafka-address localhost:9092 --input-topic accounts --output-topic accounts

## Protos

        make -C events proto
