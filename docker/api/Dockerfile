FROM golang:1.23-alpine

WORKDIR /storyden

RUN addgroup --system --gid 1001 storyden && adduser --system --uid 1001 storyden && chown storyden:storyden /storyden
USER storyden

COPY --chown=storyden:storyden . .

RUN go mod download
RUN go build -o storyden ./cmd/backend

RUN mkdir /storyden/data

VOLUME [ "/storyden/data" ]

ENTRYPOINT [ "/storyden/storyden" ]
