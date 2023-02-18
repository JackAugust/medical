cd explorer && docker-compose down -v
cd ..
docker rm -f $(docker ps -aq)
docker network prune
docker volume prune
cd fixtures && docker-compose up -d
cd ..
cd explorer && docker-compose up -d
cd ..
rm medical
go build
./medical
