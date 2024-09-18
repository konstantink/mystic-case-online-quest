FROM --platform=linux/amd64 golang:1.22-alpine3.19 AS app_builder

WORKDIR /code

ENV PORT=8080

ENV CGO_ENABLED=1
ENV GOPATH=/code
ENV GOCACHE=/go-build

RUN <<EOF
apk update
apk add gcc libc-dev
EOF

COPY ./go.mod ./go.sum /code/

RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download

COPY config/ /code/config
COPY *.go .

RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go build -o bin/mc-quest *.go

FROM --platform=linux/amd64 node:20.4-alpine3.17 AS styles_builder

WORKDIR /assets

COPY package.json yarn.lock tailwind.config.js /assets/
COPY templates/ /assets/templates
COPY assets/styles/ /assets/assets/styles

RUN yarn install --dev

RUN yarn build-css-prod
# RUN pwd

FROM --platform=linux/amd64 alpine:3.17 AS final

WORKDIR /app
RUN mkdir /app/static

EXPOSE 8080

COPY --from=app_builder /code/bin/mc-quest /app/
# COPY static/ ./static
COPY templates/ ./templates
COPY assets/audio/ ./static/audio
COPY assets/images/ ./static/images
COPY assets/js/ ./static/js
COPY static/fonts/ ./static/fonts
COPY hints/ ./hints
COPY --from=styles_builder /assets/static/styles.css /app/static/

RUN addgroup -S mysticcase && adduser -h /app -G mysticcase -s /bin/sh -S app
RUN chown -R app:mysticcase /app
USER app

CMD ["/app/mc-quest"]