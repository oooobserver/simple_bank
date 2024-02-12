# Build stage
FROM golang:1.22rc2-alpine3.18 AS builder
WORKDIR /app
COPY . . 
RUN go build -o main main.go
# Change the mirror
# RUN echo "http://mirrors.aliyun.com/alpine/v3.18/main/" > /etc/apk/repositories

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]