# Node (pnpm) ------------------------------------------------------------------
FROM node:22-slim AS ui
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack prepare pnpm@10.0.0 --activate && corepack enable
COPY . /usr/src/yt-dlp-webui

WORKDIR /usr/src/yt-dlp-webui/frontend

RUN rm -rf node_modules

RUN pnpm install
RUN pnpm run build
# -----------------------------------------------------------------------------

# Go --------------------------------------------------------------------------
FROM golang AS build

WORKDIR /usr/src/yt-dlp-webui

COPY . .
COPY --from=ui /usr/src/yt-dlp-webui/frontend /usr/src/yt-dlp-webui/frontend

RUN CGO_ENABLED=0 GOOS=linux go build -o yt-dlp-webui
# -----------------------------------------------------------------------------

# Runtime ---------------------------------------------------------------------
FROM alpine

RUN apk update && \
apk add ffmpeg ca-certificates curl wget gnutls && \
RELEASE=$(curl --silent "https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/') && \
if [ $(apk --print-arch) == "aarch64" ]; then curl -L -o yt-dlp "https://github.com/yt-dlp/yt-dlp/releases/download/$RELEASE/yt-dlp_linux_aarch64"; else curl -L -o yt-dlp "https://github.com/yt-dlp/yt-dlp/releases/download/$RELEASE/yt-dlp_linux"; fi; && \
chmod +x ./yt-dlp && \
mv ./yt-dlp /usr/bin/yt-dlp

VOLUME /downloads /config

WORKDIR /app

COPY --from=build /usr/src/yt-dlp-webui/yt-dlp-webui /app

ENV JWT_SECRET=secret

EXPOSE 3033
ENTRYPOINT [ "./yt-dlp-webui" , "--out", "/downloads", "--conf", "/config/config.yml", "--db", "/config/local.db" ]