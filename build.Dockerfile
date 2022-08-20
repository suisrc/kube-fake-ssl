FROM golang:1.18-bullseye as build
WORKDIR /build
COPY . ./
RUN go build -ldflags "-w -s" -o ./bin/runner .

FROM debian:bullseye-slim

LABEL maintainer="suisrc@outlook.com"

COPY --from=build /build/bin/runner /www/
WORKDIR /www
EXPOSE  80
ENTRYPOINT ["./runner"]
