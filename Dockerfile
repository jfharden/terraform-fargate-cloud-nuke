FROM alpine:3

LABEL Maintainer=jfharden@gmail.com
LABEL Name=fargate-cloud-nuke
LABEL Version=0.1

ARG CLOUD_NUKE_VERSION=v0.1.17
ARG CLOUD_NUKE_BINARY=cloud-nuke_linux_amd64

RUN \
  addgroup -S cloud-nuke && \
  adduser -S -G cloud-nuke cloud-nuke && \
  mkdir /cloud-nuke && \
  chown cloud-nuke:cloud-nuke /cloud-nuke

WORKDIR /cloud-nuke
USER cloud-nuke

RUN \
  wget https://github.com/gruntwork-io/cloud-nuke/releases/download/$CLOUD_NUKE_VERSION/$CLOUD_NUKE_BINARY && \
  chmod u+x $CLOUD_NUKE_BINARY && \
  ln -s $CLOUD_NUKE_BINARY cloud-nuke

# The command is set to version on purpose to mean you _have_ to provide a command override to nuke your aws account.
# This is for safety since cloud nuke is so destructive
CMD ["/cloud-nuke/cloud-nuke", "--version"]
