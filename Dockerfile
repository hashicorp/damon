# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# This Dockerfile contains multiple targets.
# Use 'docker build --target=<name> .' to build one.

# ===================================
#   Non-release images.
# ===================================

# devbuild compiles the binary
# -----------------------------------
FROM golang:1.19.1 AS devbuild

# Disable CGO to make sure we build static binaries
ENV CGO_ENABLED=0

# Escape the GOPATH
WORKDIR /build
COPY . ./
RUN go build -o damon ./cmd/damon

# dev runs the binary from devbuild
# -----------------------------------
FROM alpine:3.15 AS dev

COPY --from=devbuild /build/damon /bin/
COPY ./scripts/docker-entrypoint.sh /

ENTRYPOINT ["/docker-entrypoint.sh"]


# ===================================
#   Release images.
# ===================================

FROM alpine:3.15 AS release

ARG PRODUCT_NAME=damon
ARG PRODUCT_VERSION
ARG PRODUCT_REVISION
# TARGETARCH and TARGETOS are set automatically when --platform is provided.
ARG TARGETOS TARGETARCH

LABEL maintainer="Nomad Team <nomad@hashicorp.com>"
LABEL version=${PRODUCT_VERSION}
LABEL revision=${PRODUCT_REVISION}

COPY dist/$TARGETOS/$TARGETARCH/damon /bin/
COPY ./scripts/docker-entrypoint.sh /

# Create a non-root user to run the software.
RUN addgroup $PRODUCT_NAME && \
    adduser -S -G $PRODUCT_NAME $PRODUCT_NAME

USER $PRODUCT_NAME
ENTRYPOINT ["/docker-entrypoint.sh"]

# ===================================
#   Set default target to 'dev'.
# ===================================
FROM dev
