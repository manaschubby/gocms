FROM golang:1.26-bookworm

# Install dev dependencies and Air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Speed up module downloads
COPY go.mod go.sum ./
RUN go mod download


EXPOSE 7467

# Environment: bind app to 0.0.0.0:7467 by default
ENV HTTP_ADDR=":7467"
