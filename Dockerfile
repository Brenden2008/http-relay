FROM scratch

COPY docker/passwd /etc/passwd
COPY docker/httprelay /httprelay

USER nobody

EXPOSE 8800
ENTRYPOINT ["/httprelay"]