smcp

docker build . -t smcp

docker container run -name smcp-instance -e REDIS_HOST='192.168.0.106:6379' smcp
docker container run -name smcp-instance --network=host smcp