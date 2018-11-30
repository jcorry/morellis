FROM golang:1.11

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY . /go/src/app

ARG SSH_KEY
RUN mkdir -p /root/.ssh
RUN echo "$SSH_KEY" > /root/.ssh/id_rsa
RUN chmod 400 /root/.ssh/id_rsa

RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/

RUN go get -v