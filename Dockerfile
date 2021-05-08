FROM golang:1.16-alpine
WORKDIR /go/src

COPY ./ ./
RUN go mod download
RUN apk add --no-cache gcc musl-dev