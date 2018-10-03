FROM golang:1.10 as builder

ENV GOPATH='/go/src/app/vendor:/go'\
    CGO_ENABLED='0'\
    GOOS='linux'\
    GOARCH='amd64'

WORKDIR /go/src/app
COPY . .

RUN adduser --system --no-create-home --quiet appuser &&\
    go get && \
    go build -a -installsuffix cgo -ldflags="-w -s" -o /bin/gitcount

FROM scratch

VOLUME ["/repo"]

COPY --from=builder /bin/gitcount /bin/gitcount
COPY --from=builder /etc/passwd /etc/passwd

USER appuser

CMD ["/bin/gitcount","-dir","/repo"]