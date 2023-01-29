#!/bin/bash
# Deletes the tag (branch image) passed in as an environment variable

if [ -z $TAG ]; then 
    echo "\e[32mDockerHub tag is not set...\e[0m"
    exit 1
fi

echo "> Requesting token from "
TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${DOCKERHUB_USERNAME}'", "password": "'${DOCKERHUB_PASSWORD}'"}' https://hub.docker.com/v2/users/login/ | grep -Po '"token":"*\K[^"]*')

echo "> Requesting deletion of $TAG from $REPO -> https://hub.docker.com/v2/repositories/${REPO}/tags/${TAG}/"
curl -i -X DELETE \
  -H "Accept: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  "https://hub.docker.com/v2/repositories/${REPO}/tags/${TAG}/"
