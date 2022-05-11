# This Dockerfile contains multiple targets.
# Use 'docker build --target=<name> .' to build one.

# ===================================
#   Non-release images.
# ===================================

# devbuild compiles the binary
# -----------------------------------
FROM golang:latest AS devbuild
# Escape the GOPATH
WORKDIR /build
COPY . ./
RUN go build -o damon ./cmd/damon


# dev runs the binary from devbuild
# -----------------------------------
FROM alpine:latest AS dev
RUN apk add --no-cache git libc6-compat
COPY --from=devbuild /build/damon /bin/

ENTRYPOINT ["/bin/damon"]
CMD ["-v"]


# ===================================
#   Release images.
# ===================================

FROM alpine:latest AS release

ARG PRODUCT_NAME=damon
ARG PRODUCT_VERSION
ARG PRODUCT_REVISION
# TARGETARCH and TARGETOS are set automatically when --platform is provided.
ARG TARGETOS TARGETARCH

LABEL maintainer="Nomad Team <nomad@hashicorp.com>"
LABEL version=${PRODUCT_VERSION}
LABEL revision=${PRODUCT_REVISION}

RUN apk add --no-cache git libc6-compat
COPY dist/$TARGETOS/$TARGETARCH/damon /bin/

# Create a non-root user to run the software.
RUN addgroup $PRODUCT_NAME && \
    adduser -S -G $PRODUCT_NAME $PRODUCT_NAME

USER $PRODUCT_NAME
ENTRYPOINT ["/bin/damon"]
CMD ["-v"]

# ===================================
#   Set default target to 'dev'.
# ===================================
FROM dev
