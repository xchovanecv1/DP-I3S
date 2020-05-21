#/bin/bash

#cd stack
git fetch
git checkout master
git checkout .
git pull
#docker login registry.gitlab.com -u user -p password
docker-compose -f docker-compose.yml up -d

