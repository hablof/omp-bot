# Builder

ARG GITHUB_PATH=github.com/hablof/logistic-package-api-bot

FROM golang:1.19-alpine AS builder

WORKDIR /home/${GITHUB_PATH}

RUN apk add --update make git protoc protobuf protobuf-dev curl

COPY . .
RUN make build

# telegram bot

FROM alpine:latest as bot
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/bot .
# COPY ./config.yml .

RUN chown root:root bot

CMD ["./bot"]
