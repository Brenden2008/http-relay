ARG PROJ_NS=gitlab.com/jonas.jasas
ARG PROJ_NAME=httprelay
ARG PROJ_BIN_PATH=/$PROJ_NAME

################################################################################
FROM golang:alpine
RUN apk update && apk add --no-cache git

ARG PROJ_NS
ARG PROJ_NAME
ARG PROJ_BIN_PATH

RUN go get -d $PROJ_NS/$PROJ_NAME/...
WORKDIR $GOPATH/src/$PROJ_NS/$PROJ_NAME
RUN go get -d ./cmd/...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $PROJ_BIN_PATH ./cmd/...

RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

################################################################################
FROM scratch

ARG PROJ_BIN_PATH

COPY --from=0 /etc_passwd /etc/passwd
COPY --from=0 $PROJ_BIN_PATH /entrypoint

USER nobody

EXPOSE 8800
ENTRYPOINT ["/entrypoint"]