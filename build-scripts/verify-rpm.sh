#!/bin/bash

BASE_URL="https://artifactory.hashicorp.engineering/artifactory/hashicorp-crt-staging-local"

curl -u "$ARTIFACTORY_USER":"$ARTIFACTORY_TOKEN" \
    -X GET "${BASE_URL}/${REPO_NAME}/${VERSION}/${GIT_SHA}/${REPO_NAME}-${VERSION}-1.x86_64.rpm" \
    --output "${REPO_NAME}-${VERSION}-1.x86_64.rpm"

yum -y localinstall $REPO_NAME-$VERSION-1.x86_64.rpm

$REPO_NAME -v