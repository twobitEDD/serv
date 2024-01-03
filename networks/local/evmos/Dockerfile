FROM golang:1.20-bullseye as build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/twobitedd/serv

# Add source files
COPY . .

# build Ethermint
RUN make build-linux

# Final image
FROM golang:1.20-bullseye as final

WORKDIR /

RUN apt-get update

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/twobitedd/serv/build/servd /
COPY --from=build-env /go/src/github.com/twobitedd/serv/scripts/start-docker.sh /

EXPOSE 26656 26657 1317 9090 8545 8546

# Run servd by default, omit entrypoint to ease using container with servd
ENTRYPOINT ["/bin/bash", "-c"]