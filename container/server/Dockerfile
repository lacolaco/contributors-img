FROM node:12-alpine

WORKDIR /app

# Tell Puppeteer to skip installing Chrome. We'll be using the installed package.
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
# https://github.com/puppeteer/puppeteer/blob/master/docs/troubleshooting.md#running-on-alpine
# Installs latest Chromium package.
RUN apk add --no-cache \
      chromium \
      nss \
      freetype \
      freetype-dev \
      harfbuzz \
      ca-certificates \
      ttf-freefont
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

# Setup node_modules
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

# Build apps
COPY . .
RUN yarn build:all:staging

CMD node ./dist/apps/server/main.js
