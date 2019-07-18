FROM golang:1.12
ENV GO111MODULE=on

WORKDIR /app

ADD . .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o morellis cmd/api/*.go

FROM alpine:latest

COPY --from=0 /app/morellis /app/

ENTRYPOINT ./app/morellis

EXPOSE 4001