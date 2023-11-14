FROM golang:1.18.8-alpine3.16 AS base
RUN apk add build-base

WORKDIR /src

COPY . .

RUN go mod download

EXPOSE 8000

RUN go build -o main *.go

# Deploy

FROM alpine:latest

WORKDIR /src

COPY --from=base /src/main ./

EXPOSE 8000

ENTRYPOINT [ "./main" ]