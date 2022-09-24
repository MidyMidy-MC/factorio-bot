FROM golang:latest as build

ENV CGO_ENABLED=0
COPY ./ /bot
RUN cd /bot && go build -tags netgo -trimpath -ldflags "-w -s" -o bot

FROM scratch

ENTRYPOINT [ "./bot" ]
COPY --from=build /bot/bot ./bot