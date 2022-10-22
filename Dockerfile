FROM golang:1.19-alpine AS builder

RUN apk add --no-cache \
	make=4.3-r0 \
	gcc=11.2.1_git20220219-r2 \
	musl-dev=1.2.3-r0

WORKDIR /src

COPY . .
RUN make build

FROM alpine:3.16

WORKDIR /outage

COPY --from=builder /src/cmd/outage /usr/local/bin
COPY --from=builder /src/examples/test.yml conf.yml

EXPOSE 9999
VOLUME /outage

CMD ["outage", "--bind=0.0.0.0:9999", "-pvv", "-c", "/outage/conf.yml"]
