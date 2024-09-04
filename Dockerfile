ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.0
ARG NAME=bootstrap

ARG GOCACHE=/root/.cache/go-build
ARG ASM_FLAGS="-trimpath"
ARG GC_FLAGS="-trimpath"
ARG LD_FLAGS="-w -s -extldflags '-static'"

# amd64
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder-amd64
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
COPY go.mod go.sum ./

ENV GOARCH=amd64
RUN go mod vendor
RUN --mount=type=cache,target="${GOCACHE}" env GOOS=linux GOARCH=amd64 \
      go build \
        -mod=vendor \
        -asmflags="${ASM_FLAGS}" \
        -ldflags="${LD_FLAGS}"   \
        -gcflags="${GC_FLAGS}"   \
        -o /bin/bootstrap           
RUN apk add --no-cache upx && upx --best --lzma /bin/bootstrap

FROM scratch AS amd64
COPY --from=builder-amd64 /bin/bootstrap /bootstrap
ENTRYPOINT ["/bootstrap"]

# arm64
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder-arm64
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
COPY go.mod go.sum ./

ENV GOARCH=arm64
RUN go mod vendor
RUN --mount=type=cache,target="${GOCACHE}" env GOOS=linux GOARCH=arm64 \
      go build \
        -mod=vendor \
        -asmflags="${ASM_FLAGS}" \
        -ldflags="${LD_FLAGS}"   \
        -gcflags="${GC_FLAGS}"   \
        -o /bin/bootstrap             
RUN apk add --no-cache upx && upx --best --lzma /bin/bootstrap

FROM scratch AS arm64
COPY --from=builder-arm64 /bin/bootstrap /bootstrap
ENTRYPOINT ["/bootstrap"] 
