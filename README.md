smcp

docker build . -t smcp

docker container run -name smcp-instance --restart unless-stopped -e REDIS_HOST='192.168.0.106:6379' smcp
docker container run -name smcp-instance --restart unless-stopped --network=host smcp