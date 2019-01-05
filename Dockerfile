FROM golang:latest
LABEL maintainer="suzukenz"

ENV APP_NAME discord-gce-manager
ENV USER manager
ENV HOME /home/${USER}

RUN useradd -m ${USER}
USER ${USER}

WORKDIR ${HOME}

COPY bin/linux/${APP_NAME} .
COPY entrypoint.sh .

EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]
