FROM golang:1.19-alpine AS builder

WORKDIR /src

COPY . .
RUN \
	go mod tidy && \
	go build -ldflags="-s -w" -o cmd/outage cmd/outage.go

FROM alpine:3.16
ARG BUILD_DATE
ARG GITHUB_SHA

ENV BUILD_DATE=$BUILD_DATE
ENV GITHUB_SHA=$GITHUB_SHA

EXPOSE 9999
VOLUME /outage

WORKDIR /outage

COPY --from=builder /src/cmd/outage /usr/local/bin
COPY --from=builder /src/examples/test.yml conf.yml

CMD ["outage", "--bind=0.0.0.0:9999", "-pvv", "-c", "/outage/conf.yml"]
