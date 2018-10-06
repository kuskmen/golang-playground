FROM golang:1.11

# Add a non-privileged user
RUN useradd -u 10001 appuser

RUN mkdir -p /go/src/github.com/kuskmen/golang-playground
ADD . /go/src/github.com/kuskmen/golang-playground
WORKDIR /go/src/github.com/kuskmen/golang-playground

# build the binary with go build
RUN CGO_ENABLED=0 go build -o bin/golang-playground github.com/kuskmen/golang-playground/cmd/golang-playground

FROM scratch

ENV PORT 8080
ENV DIAGNOSTICTS_PORT 8585

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=0 /etc/passwd /etc/passwd
USER appuser

COPY --from=0 /go/src/github.com/kuskmen/golang-playground/bin/golang-playground /golang-playground
EXPOSE $PORT
EXPOSE $DIAGNOSTICTS_PORT

CMD ["/golang-playground"]