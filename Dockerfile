FROM debian:buster-slim

LABEL maintainer="suisrc@outlook.com"

ADD ["bin/runner", "/www/"]
WORKDIR /www
EXPOSE  80
ENTRYPOINT ["./runner"]
