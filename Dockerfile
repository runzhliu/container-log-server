FROM golang:alpine as builder
RUN apk --no-cache add git
WORKDIR /app
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest as prod
RUN apk --no-cache add ca-certificates
COPY --from=0 /app/container-log-server .
CMD ["/container-log-server"]