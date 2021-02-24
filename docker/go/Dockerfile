FROM golang:1.15.5-buster

ARG WORK_DIR=""
ARG HTTP_PROXY=""
ARG HTTPS_PROXY=""

WORKDIR ${WORK_DIR}

ENV HTTP_PROXY=${HTTP_PROXY}
ENV HTTPS_PROXY=${HTTPS_PROXY}
ENV http_proxy=${HTTP_PROXY}
ENV https_proxy=${HTTPS_PROXY}
ENV GO111MODULE=on

RUN apt-get -y update && \
    apt-get install -y \
        make \
        zip \
        strace

COPY Makefile Makefile

RUN make install-dev && \
    make install-cobra
