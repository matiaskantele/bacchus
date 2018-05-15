FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git && \
    go get ./ && \
    go build -o bacchus

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/app/bacchus .
RUN apk update && \
    apk add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

EXPOSE 8080
ENTRYPOINT [ "./bacchus" ]
