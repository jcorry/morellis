FROM golang:1.12
ENV GO111MODULE=on

WORKDIR /app

ADD . .

COPY . .
RUN ls
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o morellis cmd/api/*.go

FROM alpine:latest

COPY --from=0 /app/morellis /app/
COPY --from=0 /app/morellis-api-a857005c832b.json /app/morellis-api-a857005c832b.json
RUN ls /app/
ENTRYPOINT ./app/morellis

EXPOSE 4001

# build it: docker build -t jcorry/morellis-api:v1 .
