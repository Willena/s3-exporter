FROM golang:1.17-alpine as build
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o s3-exporter

FROM alpine:3.15

WORKDIR /app

COPY --from=build /app/s3-exporter .

ENTRYPOINT [ "./s3-exporter" ]
