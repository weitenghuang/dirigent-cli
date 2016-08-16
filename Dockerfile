FROM golang:1.6

# ENV buildDeps="git bash curl"
ENV PATH="$PATH:$GOPATH/bin"

RUN set -x \
  # && apk add -U $buildDeps \
  && go get github.com/tools/godep \
  && curl -sSL "http://storage.googleapis.com/kubernetes-release/release/v1.3.0/bin/linux/amd64/kubectl" > /usr/bin/kubectl \
  && chmod +x /usr/bin/kubectl

WORKDIR "$GOPATH/src/github.com/weitenghuang/dirigent-cli"
COPY . "$GOPATH/src/github.com/weitenghuang/dirigent-cli"

RUN godep restore

ENTRYPOINT ["/go/src/github.com/weitenghuang/dirigent-cli/docker-entrypoint.sh"]

CMD ["go", "version"]