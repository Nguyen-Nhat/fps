FROM golang:1.19-buster as builder
WORKDIR /loyalty_file_processing

ARG CI_JOB_TOKEN
RUN git config --global url."https://gitlab-ci-token:${CI_JOB_TOKEN}@git.teko.vn".insteadOf "https://git.teko.vn"
COPY go.mod go.sum ./
ENV GOPRIVATE="git.teko.vn,*.tekoapis.com"
ENV GOSUMDB="off"
RUN go mod download
COPY ./ ./

RUN go build -o bin/server cmd/server/main.go


## Today ubuntu is using minimalized image by default, using ubuntu for better compatible than alpine
FROM ubuntu:20.04
WORKDIR /loyalty_file_processing/bin/
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /loyalty_file_processing/bin/ /loyalty_file_processing/bin/
COPY --from=builder /loyalty_file_processing/migrations/ /loyalty_file_processing/bin/migrations/
COPY --from=builder /loyalty_file_processing/resources/messages/ /loyalty_file_processing/bin/resources/messages/
COPY config.tmp.yml config.yml

EXPOSE 10080
