FROM golang:1.21 as api-builder

RUN git config --global --add safe.directory /server

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

ADD ./app ./app
ADD ./cmd ./cmd
ADD ./internal ./internal
ADD ./Taskfile.yml .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/backend

FROM node:18-alpine AS base

#
# DEPENDENCIES LAYER
#

FROM base AS deps

RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY ./web/package.json ./web/yarn.lock* ./
RUN yarn --frozen-lockfile --ignore-scripts

#
# BUILDER LAYER
#

FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY ./web .

ENV NEXT_TELEMETRY_DISABLED 1
ENV NEXT_BUILD_STANDALONE true

# NOTE: A regular Next.js build will attempt to do static generation for every
# page. But we don't want to constantly be adding force-dynamic to every page so
# this ensures no SSG runs during the Docker build.
RUN yarn next experimental-compile

#
# RUNTIME IMAGE
#

FROM base AS runner
WORKDIR /app

ENV NODE_ENV production

ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 storyden
RUN adduser --system --uid 1001 storyden
RUN chown storyden:storyden /app
USER storyden

COPY --from=builder /app/public ./public

RUN mkdir .next && mkdir data

COPY --from=api-builder --chown=storyden:storyden /server/backend ./backend
COPY --from=builder --chown=storyden:storyden /app/.next/standalone ./
COPY --from=builder --chown=storyden:storyden /app/.next/static ./.next/static


EXPOSE 3000

ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

VOLUME [ "/data" ]

# Instruct the Go process to run the Next.js server too.
ENV RUN_FRONTEND "server.js"

ENTRYPOINT ["./backend"]
