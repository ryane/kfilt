FROM golang:1.12-alpine3.9 as builder

RUN apk --no-cache add make git gcc musl-dev && rm -rf /var/cache/apk/*
WORKDIR /go/src/github.com/ryane/kfilt
RUN GO111MODULE=off go get github.com/ahmetb/govvv
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN govvv install -pkg github.com/ryane/kfilt/cmd

FROM alpine:3.9
RUN apk --no-cache add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /go/bin/kfilt /bin/kfilt
ENTRYPOINT ["/bin/kfilt"]
