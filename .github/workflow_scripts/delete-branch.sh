#!/bin/bash
# Deletes the tag (branch image) passed in as an environment variable

if [ -z $TAG ]; then 
    echo "\e[32mDockerHub tag is not set...\e[0m"
    exit 1
fi

# If the first character is a v and second character is a
# number then this is likely a release
echo "${TAG:0:1}"
echo "${TAG:1:1}"
if ! [ "${TAG:0:1}" == "v" ] || ! [[ ${TAG:1:1} =~ [0-9] ]]; then
  TAG=branch-$TAG
fi

echo "> Requesting token from https://hub.docker.com/v2/users/login/"
TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${DOCKERHUB_USERNAME}'", "password": "'${DOCKERHUB_PASSWORD}'"}' https://hub.docker.com/v2/users/login/ | grep -Po '"token":"*\K[^"]*')

echo "> Requesting deletion of $TAG from $REPO -> https://hub.docker.com/v2/repositories/${REPO}/tags/${TAG}/"
curl -i -X DELETE \
  -H "Accept: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  "https://hub.docker.com/v2/repositories/${REPO}/tags/${TAG}/"
