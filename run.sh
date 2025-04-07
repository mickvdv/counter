#!/bin/sh

docker build -t counter:latest .
docker run -p 8080:8080 --name counter counter:latest