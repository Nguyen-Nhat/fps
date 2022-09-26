FROM asia.gcr.io/teko-registry/sre/rpc-go-builder:1.14.2 as builder
WORKDIR /rpc
COPY go.mod go.sum ./
RUN go mod edit -droprequire=rpc.tekoapis.com
RUN go mod download
COPY ./ ./
RUN git clone https://git.teko.vn/shared/rpc.git shared/rpc
RUN make all

## Using ubuntu for better compatible than alpine
FROM ubuntu:20.04
WORKDIR /rpc/bin/
COPY --from=builder /rpc/bin/rpc-runtime /rpc/bin/
COPY migrations ./
EXPOSE 10080 10433
CMD ["/rpc/bin/rpc-runtime", "server"]
