FROM golang:1.16 as builder

COPY go.* /src/
WORKDIR /src
RUN go mod download

ADD . /src
RUN CGO_ENABLED=0 go build -o server cmd/main.go

# -----

FROM alpine:latest

COPY --from=builder /src/server /server

CMD ["/server"]
