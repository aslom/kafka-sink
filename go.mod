module knative.dev/eventing-contrib

go 1.14

require (
	github.com/Shopify/sarama v1.26.4
	github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2 v2.0.0
	github.com/cloudevents/sdk-go/v2 v2.1.0
	github.com/kelseyhightower/envconfig v1.4.0

)

replace (
	github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2 => github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2 v2.0.1-0.20200608152019-2ab697c8fc0b
	github.com/cloudevents/sdk-go/protocol/stan/v2 => github.com/cloudevents/sdk-go/protocol/stan/v2 v2.0.1-0.20200608152019-2ab697c8fc0b
	github.com/cloudevents/sdk-go/v2 => github.com/cloudevents/sdk-go/v2 v2.0.1-0.20200608152019-2ab697c8fc0b
)
