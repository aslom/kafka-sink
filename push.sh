source env.sh
set -e -x
docker push ${DOCKER_REPO}/kafka_sink:${KAFKA_SINK_VERSION}
