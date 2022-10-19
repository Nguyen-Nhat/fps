FROM golang:1.19-buster as builder
WORKDIR /loyalty_file_processing

COPY go.mod go.sum ./
ENV GOPRIVATE="git.teko.vn,go.tekoapis.com"
RUN go mod download
COPY ./ ./

RUN go build -o bin/server cmd/server/main.go


## Today ubuntu is using minimalized image by default, using ubuntu for better compatible than alpine
FROM ubuntu:20.04
WORKDIR /loyalty_file_processing/bin/
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /loyalty_file_processing/bin/ /loyalty_file_processing/bin/
COPY config.tmp.yml config.yml

EXPOSE 10080
