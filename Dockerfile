FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/tt-itk ./cmd

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /bin/tt-itk /app/tt-itk
COPY migrations /app/migrations

EXPOSE 8081

CMD ["/app/tt-itk"]
