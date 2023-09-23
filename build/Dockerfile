ARG DEBIAN_VERSION=bullseye

# Using arm32v7 image, as mmal is not available on 64-bit build of libraspberrypi-dev.
FROM arm32v7/golang:${DEBIAN_VERSION}

# Need to repeat this here as the ARG is cleared after a FROM.
ARG DEBIAN_VERSION

RUN apt-get -y update && \
    apt-get -y install --no-install-recommends crossbuild-essential-armhf

COPY sources.list /etc/apt/
COPY raspi.list /etc/apt/sources.list.d/
COPY raspberrypi-archive-stable.gpg /etc/apt/trusted.gpg.d/
COPY raspbian-archive-keyring.gpg /etc/apt/trusted.gpg.d/

RUN sed -i "s/RELEASE/${DEBIAN_VERSION}/g" "/etc/apt/sources.list" && \
    sed -i "s/RELEASE/${DEBIAN_VERSION}/g" "/etc/apt/sources.list.d/raspi.list"

RUN apt-get update && \
    apt-get dist-upgrade -y

RUN apt-get -y install libraspberrypi-dev

COPY gitconfig /root/.gitconfig

RUN go env -w \
    CGO_ENABLED=1 \
    CC=arm-linux-gnueabihf-gcc \
    GOARCH=arm \
    GOARM=6 \
    GOPRIVATE=github.com/thingify-app/* \
    GOCACHE=/gocache \
    GOMODCACHE=/gomodcache

WORKDIR /build
CMD go build
