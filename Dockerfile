FROM alpine:latest
COPY daoctl /
COPY daoctl.yaml /
RUN chmod +x "/daoctl.yaml"
ENTRYPOINT ["/daoctl", "serve"]