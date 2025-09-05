#
# BACKEND API LAYER
#

FROM golang:1.25-alpine AS api-builder

WORKDIR /build

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

COPY ./app ./app
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=0 GOOS=linux go build -o /build/storyden ./cmd/backend

#
# RUNTIME IMAGE
#

FROM alpine:3.22

WORKDIR /storyden

RUN addgroup --system --gid 1001 storyden && \
    adduser --system --uid 1001 storyden && \
    chown storyden:storyden /storyden

USER storyden

COPY --from=api-builder --chown=storyden:storyden /build/storyden /storyden/storyden

RUN mkdir /storyden/data

VOLUME [ "/storyden/data" ]

ENTRYPOINT [ "/storyden/storyden" ]
