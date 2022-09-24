FROM golang:latest as build

ENV CGO_ENABLED=0
COPY ./ /bot
RUN cd /bot && go build -tags netgo -trimpath -ldflags "-w -s" -o bot

FROM alpine as certs

RUN apk add -U --no-cache ca-certificates

FROM scratch

ENTRYPOINT [ "./bot" ]
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bot/bot ./bot