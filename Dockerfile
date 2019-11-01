FROM alpine/git AS vcs
RUN git clone https://gitlab.com/jonas.jasas/httprelay.git /httprelay

FROM golang AS build
RUN echo 'nobody:x:65534:65534:Nobody:/:' > /passwd
WORKDIR $GOPATH/src/httprelay
COPY --from=vcs /httprelay/ .
RUN go get ./...
RUN date --date="`git show -s --format=%ci`" "+%Y%m%d-%H%M" > version
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.Version=`cat version`" -o /httprelay ./cmd/...

FROM scratch
COPY --from=build /passwd /etc/passwd
COPY --from=build /httprelay .
USER nobody
EXPOSE 8080
ENTRYPOINT ["/httprelay"]