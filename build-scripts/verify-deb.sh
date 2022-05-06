#!/bin/bash

apt-get update && \
apt-get install -y curl gnupg2 lsb-release software-properties-common && \
apt-get clean all

export BASE_URL="https://artifactory.hashicorp.engineering/artifactory/hashicorp-crt-staging-local"

curl -u "$ARTIFACTORY_USER":"$ARTIFACTORY_TOKEN" \
    -X GET "${BASE_URL}/${REPO_NAME}/${VERSION}/${GIT_SHA}/${REPO_NAME}_${VERSION}-1_amd64.deb" \
    --output "${REPO_NAME}-${VERSION}-1.amd64.deb"

apt-get -y install ./$REPO_NAME-$VERSION-1.amd64.deb

$REPO_NAME -v