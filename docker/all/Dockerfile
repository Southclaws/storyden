FROM golang:latest as api-builder

RUN git config --global --add safe.directory /server

WORKDIR /server
ADD . .

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN task backend

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

RUN yarn build

#
# RUNTIME IMAGE
#

FROM base AS runner
WORKDIR /app

ENV NODE_ENV production

ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public

RUN mkdir .next
RUN chown nextjs:nodejs .next

COPY --from=api-builder --chown=nextjs:nodejs /server/backend.exe ./backend.exe

# Not enabled currently.
# COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

# TODO: Spin up the backend server here too.
CMD ["node", "server.js"]
