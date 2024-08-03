FROM golang:1.22.5-alpine AS base

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM scratch

COPY --from=base /app/main /

EXPOSE 3031

CMD ["/main"]