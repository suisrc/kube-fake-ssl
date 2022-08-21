FROM golang:1.18-bullseye as build
WORKDIR /build
COPY . ./
RUN go build -ldflags "-w -s" -o ./bin/kube-fake-ssl .

FROM debian:bullseye-slim

LABEL maintainer="suisrc@outlook.com"

COPY --from=build /build/bin/kube-fake-ssl /www/
WORKDIR /www
EXPOSE  80
ENTRYPOINT ["./kube-fake-ssl"]
