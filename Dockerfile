FROM scratch
COPY daoctl /
ENTRYPOINT ["/daoctl", "serve"]