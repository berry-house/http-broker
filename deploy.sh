#/bin/bash

echo "Compiling code..."
go build || exit

DATE=$(date '+%Y-%m-%d')
IMG_NAME=berryhouse/httpbroker

echo "Removing previous \"latest\" tag..."
docker rmi $IMG_NAME:latest

echo "Building image..."
docker build -t $IMG_NAME:$DATE -t $IMG_NAME:latest .

echo "Pushing image to Docker Hub..."
docker push $IMG_NAME:$DATE

