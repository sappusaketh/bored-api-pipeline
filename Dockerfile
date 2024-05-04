ARG GOVERSION=1.22.2
FROM --platform=$BUILDPLATFORM golang:${GOVERSION} AS build

WORKDIR /app

ARG GOPROXY

ENV GOPROXY=${GOPROXY}
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH make build

FROM golang:${GOVERSION}-alpine3.19

COPY --from=build /app/bored-api-pipeline /usr/bin/bored-api-pipeline
COPY --from=build /app/config ./config

ENTRYPOINT ["bored-api-pipeline"]
