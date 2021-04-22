FROM golang:1.16-alpine as builder
RUN apk add make binutils

COPY / /work
WORKDIR /work
RUN make

FROM alpine:3.13
COPY --from=builder /work/bin/partition-watchdog /partition-watchdog
USER 65534
ENTRYPOINT ["/partition-watchdog","check"]
