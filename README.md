# Kafka Sink: HTTP Sink to Apache Kafka 

Th Kafka Sink listens on HTTP port for CloudEvents that are then sent 
to the topic in Apache Kafka broker.


## Quick testing without Kubernetes

RUn Kafka Sink from command line:

```
export KAFKA_SERVER=localhost:9092
export KAFKA_TOPIC=test-topic
go run cmd/kafka_sink/main.go
```

Send CloudEvent to Kafka Sink:

```
curl -v "http://localhost:8080" \
  -X POST \
  -H "Ce-Id: test-kafka-sink" \
  -H "Ce-Specversion: 1.0" \
  -H "Ce-Type: greeting" \
  -H "Ce-Source: not-sendoff" \
  -H "Content-Type: application/json" \
  -d '{"msg":"Hello Kafka Sink!"}'
```


## Deploying Kafka Sink as Kubernetes Service

Edit deployment YAML in service-config to point to your Kafka broker running in your cluster and Kafka topic.

Then deploy it:

```
kubectl apply -f service-config/
```

Check that it is running - expected output:

```
kubectl get pods
NAME                          READY   STATUS    RESTARTS   AGE
kafka-sink-694fd8f689-smwq4   1/1     Running   0          8s
```

```
kubectl logs kafka-sink-694fd8f689-smwq4
2020/07/22 22:13:03 Using HTTP PORT=8080
2020/07/22 22:13:03 Sinking to KAFKA_SERVER=my-cluster-kafka-bootstrap.kafka:9092 KAFKA_TOPIC=knative-sink-topic
2020/07/22 22:13:03 Ready to receive
```

To test follow the same approach as sending events to broker: https://knative.dev/docs/eventing/getting-started/#sending-events-to-the-broker

Create CLI pod:

```
kubectl apply --filename - << END
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: curl
  name: curl
spec:
  containers:
    # This could be any image that we can attach into and has curl.
  - image: radial/busyboxplus:curl
    imagePullPolicy: IfNotPresent
    name: curl
    resources: {}
    stdin: true
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    tty: true
END
```

Get into it:

```
kubectl attach curl -it
Defaulting container name to curl.
Use 'kubectl describe pod/curl -n default' to see all of the containers in this pod.
If you don't see a command prompt, try pressing enter.
[ root@curl:/ ]$
```

Send (modify URL to replace default with current namespace)

```
curl -v "http://kafka-sink.default.svc.cluster.local" \
  -X POST \
  -H "Ce-Id: kafka-sink-event-1" \
  -H "Ce-Specversion: 1.0" \
  -H "Ce-Type: greeting" \
  -H "Ce-Source: not-sendoff" \
  -H "Content-Type: application/json" \
  -d '{"msg":"Hello Knative Kafka Sink!"}'
```

Check Kafka Sink logs:

```
k logs -f kafka-sink-694fd8f689-smwq4
2020/07/22 22:13:03 Using HTTP PORT=8080
2020/07/22 22:13:03 Sinking to KAFKA_SERVER=my-cluster-kafka-bootstrap.kafka:9092 KAFKA_TOPIC=knative-sink-topic
2020/07/22 22:13:03 Ready to receive
2020/07/22 22:19:46 Received message
2020/07/22 22:19:46 Sending message to Kafka
2020/07/22 22:19:47 Ready to receive
```

And verify that event was received by Kafka, for example

```
kubectl -n kafka run kafka-consumer -ti --image=strimzi/kafka:0.17.0-kafka-2.4.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic knative-sink-topic --from-beginning
If you don't see a command prompt, try pressing enter.
{"msg":"Hello Knative Kafka Sink!"}
```

## Deploying Kafka Sink as Kubernetes Service


