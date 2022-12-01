FROM golang:alpine as builder

LABEL maintainer="insomniadev <insomniadevlabs@gmail.com>"

COPY . $GOPATH/src/github.com/insomniadev/martian
WORKDIR $GOPATH/src/github.com/insomniadev/martian

RUN apk add git
RUN go build ./cmd/martian
RUN ls -al
RUN pwd

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/insomniadev/martian/martian /martian

RUN mkdir config
RUN pwd && ls -al
WORKDIR /

CMD ["./martian"]
EXPOSE 9000