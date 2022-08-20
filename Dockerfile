FROM debian:buster-slim

#ARG http://mirrors.aliyun.com http://mirrors.163.com
ARG LINUX_MIRRORS=http://mrs.dev.sims-cn.com
# set version label
LABEL maintainer="suisrc@outlook.com"

# update linux
RUN if [ ! -z ${LINUX_MIRRORS+x} ]; then \
        mv /etc/apt/sources.list /etc/apt/sources.list.bak && \
        echo "deb ${LINUX_MIRRORS}/debian/ buster main non-free contrib" >>/etc/apt/sources.list &&\
        echo "deb ${LINUX_MIRRORS}/debian/ buster-updates main non-free contrib" >>/etc/apt/sources.list &&\
        echo "deb ${LINUX_MIRRORS}/debian/ buster-backports main non-free contrib" >>/etc/apt/sources.list &&\
        echo "deb-src ${LINUX_MIRRORS}/debian/ buster main non-free contrib" >>/etc/apt/sources.list &&\
        echo "deb-src ${LINUX_MIRRORS}/debian/ buster-updates main non-free contrib" >>/etc/apt/sources.list &&\
        echo "deb-src ${LINUX_MIRRORS}/debian/ buster-backports main non-free contrib" >>/etc/apt/sources.list; \
    fi &&\
    apt-get -o Acquire::Check-Valid-Until=false update && apt-get install --no-install-recommends -y ca-certificates &&\
    rm -rf /tmp/* /var/tmp/* /var/lib/apt/lists/* && mkdir -p /www/

ADD ["bin/runner", "/www/"]
WORKDIR /www
EXPOSE  80
ENTRYPOINT ["./runner"]
