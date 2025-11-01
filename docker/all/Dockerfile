#
# BACKEND API LAYER
#

FROM docker.io/golang:1.25-alpine AS api-builder

WORKDIR /server

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download

COPY ./app ./app
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    CGO_ENABLED=0 GOOS=linux go build -o /server/storyden-api ./cmd/backend

FROM docker.io/node:22.21.1-alpine AS base

#
# DEPENDENCIES LAYER
#

FROM base AS deps

RUN apk add --no-cache libc6-compat
WORKDIR /storyden

COPY ./web/package.json ./web/yarn.lock* ./web/.yarnrc.yml ./
RUN --mount=type=cache,target=/root/.yarn \
    corepack enable && \
    YARN_CACHE_FOLDER=/root/.yarn yarn install --immutable --network-timeout 1000000

#
# BUILDER LAYER
#

FROM base AS builder
WORKDIR /storyden
COPY --from=deps /storyden/node_modules ./node_modules
COPY ./web .

ENV NEXT_TELEMETRY_DISABLED=1
ENV NEXT_BUILD_STANDALONE=true

# NOTE: A regular Next.js build will attempt to do static generation for every
# page. But we don't want to constantly be adding force-dynamic to every page so
# this ensures no SSG runs during the Docker build.
RUN --mount=type=cache,target=/storyden/.next/cache \
    corepack enable && yarn next build --experimental-build-mode=compile

#
# RUNTIME IMAGE
#

FROM base AS runner
WORKDIR /storyden

ENV NODE_ENV=production

ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 storyden && \
    adduser --system --uid 1001 storyden && \
    chown storyden:storyden /storyden
USER storyden

COPY --from=builder /storyden/public ./public

RUN mkdir .next && mkdir data

COPY --from=api-builder --chown=storyden:storyden /server/storyden-api ./storyden-api
COPY --from=builder --chown=storyden:storyden /storyden/.next/standalone ./
COPY --from=builder --chown=storyden:storyden /storyden/.next/static ./.next/static

EXPOSE 8000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

VOLUME [ "/storyden/data" ]

# Instruct the API process to run the Next.js server and proxy requests to it.
ENV RUN_FRONTEND="server.js"
ENV PROXY_FRONTEND_ADDRESS="http://localhost:3000"

# The fullstack image uses the backend as a proxy for the frontend, this means
# that the address for both are the same. Change via -e if deployed publicly.
ENV PUBLIC_API_ADDRESS="http://localhost:8000"
ENV PUBLIC_WEB_ADDRESS="http://localhost:8000"

# You won't need to change this ever, if you do then this image is not for you.
# This tells the Next.js server to make local HTTP requests directly to the API
# server instead of going via some external network and through the API proxy.
# This setting is pretty much only used when both services run side by side.
ENV SSR_API_ADDRESS="http://localhost:8000"

ENTRYPOINT ["/storyden/storyden-api"]
