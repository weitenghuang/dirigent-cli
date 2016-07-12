FROM golang:1.6-alpine

ENV buildDeps="git bash"
ENV PATH="$PATH:$GOPATH/bin"

RUN set -x \
  && apk add -U $buildDeps \
  && go get github.com/tools/godep

WORKDIR "$GOPATH/src/github.com/weitenghuang/dirigent-cli"
COPY . "$GOPATH/src/github.com/weitenghuang/dirigent-cli"

ENTRYPOINT ["/go/src/github.com/weitenghuang/dirigent-cli/docker-entrypoint.sh"]

CMD ["go", "version"]