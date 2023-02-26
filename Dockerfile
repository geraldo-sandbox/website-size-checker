FROM golang:1.20.1-alpine3.17 AS builder

RUN apk update && apk add --no-cache make bash gcc musl-dev libc-dev ca-certificates curl build-base
RUN adduser -D -g '' appuser
WORKDIR /app
COPY . .
RUN make build


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/build/bin/wsc /app/build/bin/wsc

USER appuser

ENTRYPOINT ["/app/build/bin/wsc"]

