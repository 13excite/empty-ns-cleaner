FROM plugins/base:linux-amd64


COPY ns-cleaner /bin/

ENTRYPOINT ["/bin/ns-cleaner"]
