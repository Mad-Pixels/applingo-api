ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.0
ARG NAME=bootstrap

ARG GOCACHE=/root/.cache/go-build
ARG ASM_FLAGS="-trimpath"
ARG GC_FLAGS="-trimpath"
ARG LD_FLAGS="-w -s -extldflags '-static'"

# amd64
FROM --platform=linux/amd64 golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder-amd64
ARG FUNC_NAME
ARG ASM_FLAGS
ARG GC_FLAGS
ARG LD_FLAGS
ARG GOCACHE
ARG NAME

WORKDIR /go/src/${NAME}
COPY ./cmd/${FUNC_NAME}/ ./
COPY ./vendor ./vendor
COPY ./pkg ./pkg
COPY ./internal ./internal
COPY ./data ./data
COPY go.mod go.sum ./

ENV GOARCH=amd64
RUN go mod vendor
RUN --mount=type=cache,target="${GOCACHE}" env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
      go build \
        -mod=vendor \
        -asmflags="${ASM_FLAGS}" \
        -ldflags="${LD_FLAGS}"   \
        -gcflags="${GC_FLAGS}"   \
        -o /bin/bootstrap           
RUN apk add --no-cache upx && upx --best --lzma /bin/bootstrap
RUN wget -O /tmp/aws-ca-bundle.pem https://curl.se/ca/cacert.pem

FROM scratch AS amd64
COPY --from=builder-amd64 /bin/bootstrap /bootstrap
COPY --from=builder-amd64 /tmp/aws-ca-bundle.pem /etc/ssl/certs/aws-ca-bundle.pem
ENTRYPOINT ["/bootstrap"]

# arm64
FROM --platform=linux/arm64/v8 golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder-arm64
ARG FUNC_NAME
ARG ASM_FLAGS
ARG GC_FLAGS
ARG LD_FLAGS
ARG GOCACHE
ARG NAME

WORKDIR /go/src/${NAME}
COPY ./cmd/${FUNC_NAME}/ ./
COPY ./vendor ./vendor
COPY ./pkg ./pkg
COPY ./internal ./internal
COPY ./data ./data
COPY go.mod go.sum ./

ENV GOARCH=arm64
RUN go mod vendor
RUN --mount=type=cache,target="${GOCACHE}" env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
      go build \
        -mod=vendor \
        -asmflags="${ASM_FLAGS}" \
        -ldflags="${LD_FLAGS}"   \
        -gcflags="${GC_FLAGS}"   \
        -o /bin/bootstrap             
RUN apk add --no-cache upx && upx --best --lzma /bin/bootstrap
RUN wget -O /tmp/aws-ca-bundle.pem https://curl.se/ca/cacert.pem

FROM scratch AS arm64
COPY --from=builder-arm64 /bin/bootstrap /bootstrap
COPY --from=builder-arm64 /tmp/aws-ca-bundle.pem /etc/ssl/certs/aws-ca-bundle.pem
ENTRYPOINT ["/bootstrap"]