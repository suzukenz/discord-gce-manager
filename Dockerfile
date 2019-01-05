FROM golang:latest
LABEL maintainer="suzukenz"

ENV USER manager
ENV HOME /home/${USER}

RUN useradd -m ${USER}
USER ${USER}

WORKDIR ${HOME}

COPY bin/linux/discord-bot .

ENTRYPOINT [ "./discord-bot" ]
