package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudeventskafka "github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudeventshttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int `envconfig:"PORT" default:"8080"`

	// KafkaServer URL to connect to the Kafka server.
	KafkaServer string `envconfig:"KAFKA_SERVER" default:"127.0.0.1:9092" required:"true"`

	// Subject is the nats subject to publish cloudevents on.
	Topic string `envconfig:"KAFKA_TOPIC" default:"test-topic" required:"true"`
}

const (
	count = 10
)

func main() {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	//[]string{"127.0.0.1:9092"}
	sender, err := kafka_sarama.NewSender([]string{env.KafkaServer}, saramaConfig, "test-topic2")
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	defer sender.Close(context.Background())

	// c, err := cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	// if err != nil {
	// 	log.Fatalf("failed to create client, %v", err)
	// }

	// for i := 0; i < count; i++ {
	// 	e := cloudevents.NewEvent()
	// 	e.SetType("com.cloudevents.sample.sent")
	// 	e.SetSource("https://github.com/cloudevents/sdk-go/v2/samples/kafka/sender")
	// 	_ = e.SetData(cloudevents.ApplicationJSON, map[string]interface{}{
	// 		"id":      i,
	// 		"message": "Hello, World!",
	// 	})

	// 	if result := c.Send(context.Background(), e); !cloudevents.IsACK(result) {
	// 		log.Printf("failed to send: %v", err)
	// 	} else {
	// 		log.Printf("sent: %d, accepted: %t", i, cloudevents.IsACK(result))
	// 	}
	// }
	ctx := context.Background()

	// natsProtocol, err := cloudeventsnats.NewSender(env.NATSServer, env.Subject, cloudeventsnats.NatsOptions())
	// if err != nil {
	// 	log.Fatalf("failed to create nats protcol, %s", err.Error())
	// }

	// defer natsProtocol.Close(ctx)

	kafkaProtocol, err := cloudeventskafka.NewSender([]string{env.KafkaServer}, saramaConfig, env.Topic)
	// if err != nil {
	// 	log.Fatalf("failed to create nats protcol, %s", err.Error())
	// }

	// defer natsProtocol.Close(ctx)

	httpProtocol, err := cloudeventshttp.New(cloudeventshttp.WithPort(env.Port))
	if err != nil {
		log.Fatalf("failed to create http protocol: %s", err.Error())
	}

	// Pipe all messages incoming to the httpProtocol to the kafkaProtocol
	go func() {
		for {
			log.Printf("Ready to receive")
			// Blocking call to wait for new messages from httpProtocol
			message, err := httpProtocol.Receive(ctx)
			log.Printf("Received")
			if err != nil {
				if err == io.EOF {
					return // Context closed and/or receiver closed
				}
				log.Printf("Error while receiving a message: %s", err.Error())
			}
			// Send message directly to natsProtocol
			// err = natsProtocol.Send(ctx, message)
			// if err != nil {
			// 	log.Printf("Error while forwarding the message: %s", err.Error())
			// }
			err = kafkaProtocol.Send(ctx, message)
			if err != nil {
				log.Printf("Error while forwarding the message: %s", err.Error())
			}
		}
	}()

	// Start the HTTP Server invoking OpenInbound()
	go func() {
		if err := httpProtocol.OpenInbound(ctx); err != nil {
			log.Printf("failed to StartHTTPReceiver, %v", err)
		}
	}()

	<-ctx.Done()
}
