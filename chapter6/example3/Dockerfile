FROM ubuntu:18.04
RUN export DEBIAN_FRONTEND=noninteractive DEBCONF_NONINTERACTIVE_SEEN=true && \
    apt-get update && apt-get install -y \
        wget \
        curl \
        clang \
        libpython-dev \
        libblocksruntime-dev \
        binutils \
        git \
        libc6-dev \
        libcurl4 \
        libedit2 \
        libgcc-5-dev \
        libpython2.7 \
        libsqlite3-0 \
        libstdc++-5-dev \
        libxml2 \
        pkg-config \
        tzdata \
        zlib1g-dev && \
    wget -q -O - https://git.io/vQhTU | bash && \
    wget https://storage.googleapis.com/swift-tensorflow-artifacts/releases/v0.11/rc2/swift-tensorflow-RELEASE-0.11-ubuntu18.04.tar.gz && \
    tar xzf swift-tensorflow-RELEASE-0.11-ubuntu18.04.tar.gz && \
    rm -rf swift-tensorflow-RELEASE-0.11-ubuntu18.04.tar.gz && \
    mkdir -p code

ENV GOROOT=/root/.go \
    GOPATH=/root/go \
    PATH=/root/.go/bin:/root/go/bin:$PATH
COPY golang/ swift/ code/

RUN swiftc -O -Ounchecked -enforce-exclusivity=none -emit-library code/TWAI.swift && mv libTWAI.so /usr/lib/ && cd code && go build .
ENTRYPOINT ["/code/example"]