# Build stage
FROM golang:1.22rc2-alpine3.18 AS builder
WORKDIR /app
COPY . . 
RUN go build -o main main.go
# Change the mirror
# RUN echo "http://mirrors.aliyun.com/alpine/v3.18/main/" > /etc/apk/repositories
RUN apk --no--cache add curl
# Install the go-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]