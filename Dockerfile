from golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main ./cmd/auth-service/main.go

WORKDIR /dist
RUN cp /build/main .
RUN cp /build/config.json .

EXPOSE 8886
ENTRYPOINT ["/dist/main"]
CMD ["--config=/dist/config.json", "--loglevel=debug"]