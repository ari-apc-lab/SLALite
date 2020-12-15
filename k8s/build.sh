#!/bin/bash

export IMAGE_NAME=slalite
docker build --rm=true -t $IMAGE_NAME /home/linux/gocode_2/src/SLALite
docker tag $IMAGE_NAME registry.test.euxdat.eu/euxdat/slalite
docker push registry.test.euxdat.eu/euxdat/slalite

