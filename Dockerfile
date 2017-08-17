FROM alpine:3.6

COPY randfault /usr/local/bin/randfault
ENTRYPOINT ["/usr/local/bin/randfault"]
