cd explorer && docker-compose down -v
cd ..
docker rm -f $(docker ps -aq)
docker network prune
docker volume prune
rm medical
