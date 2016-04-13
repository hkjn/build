# TODO(hkjn): Could start from scrath and add binary.
# TODO(hkjn): Should find a way to handle matrix builds for different
# base images (can be used for different CPU archs).
FROM hkjn/armv7l-golang

MAINTAINER Henrik Jonsson <me@hkjn.me>

LABEL type=infra
ENV DOMAIN build.hkjn.me
ENV BUILD_CERT /etc/letsencrypt/live/$DOMAIN/fullchain.pem
ENV BUILD_KEY /etc/letsencrypt/live/$DOMAIN/privkey.pem
ENV PORT 4430
ENV BUILD_ADDR :$PORT
# TODO(hkjn): Would be nice to split BUILD_PATHS line.
ENV BUILD_PATHS "/hkjn.me/bitcoin/b29jZWlYYWZvaGdoYWljN,/hkjn.me/build/GVwaG9oMG9vUGhvb3NhCg"
ENV BUILD_TASKS /etc/build/tasks/

USER root
RUN pacman -Syyu && \
    pacman -S --noconfirm letsencrypt

RUN letsencrypt --register-unsafely-without-email \
                --agree-tos --non-interactive --tls-sni-01-port $PORT \
                certonly $DOMAIN

RUN chown -R go:go /etc/letsencrypt
USER go
WORKDIR /go/src/hkjn.me/build/
COPY *.go ./
RUN go test -race && go vet && go build

EXPOSE $PORT
CMD ["build"]



