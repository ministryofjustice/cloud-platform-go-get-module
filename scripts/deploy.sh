#!/usr/bin/env bash

set -ex

VERSION=$1 

docker buildx build --platform linux/amd64 --no-cache -t jaskaransarkaria/cloud-platform-go-get-module:$VERSION .
docker push jaskaransarkaria/cloud-platform-go-get-module:$VERSION
gsed -i "s/VERSION/$VERSION/" deploy/app.yaml
kubectl apply -f deploy/app.yaml
gsed -i "s/$VERSION/VERSION/" deploy/app.yaml

exit 0
