source env.sh
set -e -x
docker build -t ${DOCKER_REPO}/kafka_sink:${KAFKA_SINK_VERSION} .
