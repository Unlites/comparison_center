FROM golang:1.21-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /migrate ./cmd/migrate/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /comparison_center ./cmd/comparison_center/main.go

FROM scratch

COPY --from=builder migrate /bin/migrate
COPY --from=builder comparison_center /bin/comparison_center

CMD ["/bin/comparison_center"]