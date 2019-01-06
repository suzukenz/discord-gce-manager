FROM golang:latest
LABEL maintainer="suzukenz"

ENV HOME /root

RUN apt-get update && apt-get install -y \
  busybox-static \
  && apt-get clean

WORKDIR ${HOME}

COPY bin/linux/checker .
COPY checker/crontab /var/spool/cron/crontabs/root
COPY configs.env .

ENTRYPOINT [ "busybox", "crond", "-f", "-L", "/dev/stderr" ]
