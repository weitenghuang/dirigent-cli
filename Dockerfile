FROM golang:1.6

ENV PATH="$PATH:/go/bin"

RUN set -x \
  && go get github.com/tools/godep \
  && curl -sSL "http://storage.googleapis.com/kubernetes-release/release/v1.3.0/bin/linux/amd64/kubectl" > /usr/bin/kubectl \
  && chmod +x /usr/bin/kubectl

WORKDIR "/go/src/github.com/weitenghuang/dirigent-cli"
COPY . "$GOPATH/src/github.com/weitenghuang/dirigent-cli"

# RUN godep restore
RUN CGO_ENABLED=0 GOOS=linux go build -i -x -o /go/bin/dirigent github.com/weitenghuang/dirigent-cli

ENTRYPOINT ["/go/src/github.com/weitenghuang/dirigent-cli/docker-entrypoint.sh"]

CMD ["tar", "-cf", "-", "-C", "/go/src/github.com/weitenghuang/dirigent-cli", "Dockerfile.install", "-C", "/go/src/github.com/weitenghuang/dirigent-cli", "docker-entrypoint.sh", "-C", "/usr/bin", "kubectl", "-C", "/go/bin", "dirigent"]