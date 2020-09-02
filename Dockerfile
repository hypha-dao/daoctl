FROM alpine:latest
COPY daoctl /
ENTRYPOINT ["/daoctl", "serve"]