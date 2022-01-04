FROM golang:alpine as builder

LABEL maintainer="Adam Martin <adam9190@gmail.com>"

COPY . $GOPATH/src/github.com/insomniadev/martian
WORKDIR $GOPATH/src/github.com/insomniadev/martian

RUN apk add git
# RUN go get -u github.com/golang/dep/cmd/dep;export GOOS=linux && export CGO_ENABLED=0; dep ensure
RUN go build .
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
# VOLUME /config