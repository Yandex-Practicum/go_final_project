FROM ubuntu:latest
LABEL authors="hunt"

ENTRYPOINT ["top", "-b"]