FROM golang:1.18 as builder

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod .
#COPY go.sum .
RUN go mod download

COPY cmd/ ./cmd/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main ./cmd/main.go


#--------------------------------------------------------------execution stage--------------------------------------------------------------
FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
