FROM docker.arvancloud.ir/golang AS builder

WORKDIR /backend

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD bash -c "go run ./cmd/main.go"