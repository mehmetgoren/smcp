smcp

docker build . -t smcp

docker container run --name smcp-service --restart unless-stopped -e REDIS_HOST='192.168.0.106:6379' -v /home/gokalp/Pictures/detected/:/go/src/smcp/images/ smcp
docker container run --name smcp-service --restart unless-stopped --network=host -v /home/gokalp/Pictures/detected/:/go/src/smcp/images/ smcp