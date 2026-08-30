FROM golang:1.25.5-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/exe ./cmd/app

FROM alpine:3.23
RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
WORKDIR /app
COPY --from=builder /app/exe /app/exe
COPY --from=builder /app/frontend/dist /app/frontend/dist

ENV PROJECT_ROOT=/app
EXPOSE 5050
CMD ["/app/exe"]