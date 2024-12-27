FROM golang:1.23 AS api-builder

RUN git config --global --add safe.directory /server

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

ADD ./app ./app
ADD ./cmd ./cmd
ADD ./internal ./internal
ADD ./Taskfile.yml .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/backend

FROM node:22-alpine AS base

#
# DEPENDENCIES LAYER
#

FROM base AS deps

RUN apk add --no-cache libc6-compat
WORKDIR /storyden

COPY ./web/package.json ./web/yarn.lock* ./web/.yarnrc.yml ./
RUN corepack enable && yarn install --immutable --network-timeout 1000000

#
# BUILDER LAYER
#

FROM base AS builder
WORKDIR /storyden
COPY --from=deps /storyden/node_modules ./node_modules
COPY ./web .

ENV NEXT_TELEMETRY_DISABLED 1
ENV NEXT_BUILD_STANDALONE true

# NOTE: A regular Next.js build will attempt to do static generation for every
# page. But we don't want to constantly be adding force-dynamic to every page so
# this ensures no SSG runs during the Docker build.
RUN corepack enable && yarn next build --experimental-build-mode=compile

#
# RUNTIME IMAGE
#

FROM base AS runner
WORKDIR /storyden

ENV NODE_ENV production

ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 storyden
RUN adduser --system --uid 1001 storyden
RUN chown storyden:storyden /storyden
USER storyden

COPY --from=builder /storyden/public ./public

RUN mkdir .next && mkdir data

COPY --from=api-builder --chown=storyden:storyden /server/backend ./backend
COPY --from=builder --chown=storyden:storyden /storyden/.next/standalone ./
COPY --from=builder --chown=storyden:storyden /storyden/.next/static ./.next/static

EXPOSE 3000

ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

VOLUME [ "/data" ]

# Instruct the API process to run the Next.js server and proxy requests to it.
ENV RUN_FRONTEND "server.js"
ENV PROXY_FRONTEND_ADDRESS "http://localhost:3000"

ENTRYPOINT ["./backend"]
