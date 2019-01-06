FROM golang:latest
LABEL maintainer="suzukenz"

ENV USER discordbot
ENV HOME /home/${USER}

RUN useradd -m ${USER}
USER ${USER}

WORKDIR ${HOME}

COPY bin/linux/discord-bot .
COPY configs.env .

ENTRYPOINT [ "./discord-bot" ]
