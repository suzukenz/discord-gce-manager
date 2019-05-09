FROM golang:latest
LABEL maintainer="suzukenz"

ENV USER scheduled-checker
ENV HOME /home/${USER}

RUN useradd -m ${USER}
USER ${USER}

WORKDIR ${HOME}

COPY bin/scheduled-checker .

ENTRYPOINT [ "./scheduled-checker" ]
