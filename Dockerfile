# FROM golang:1.12-alpine3.9 as builder

# RUN apk --no-cache add make git gcc musl-dev && rm -rf /var/cache/apk/*
# WORKDIR /go/src/github.com/ryane/kfilt
# ENV GO111MODULE=on
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# COPY . .
# RUN go install

# FROM alpine:3.9
# RUN apk --no-cache add ca-certificates && rm -rf /var/cache/apk/*
# COPY --from=builder /go/bin/kfilt /bin/kfilt
# ENTRYPOINT ["/bin/kfilt"]

FROM scratch
COPY kfilt /
ENTRYPOINT ["/kfilt"]
