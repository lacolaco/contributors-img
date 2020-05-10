FROM node:12-alpine as install

# Tell Puppeteer to skip installing Chrome. We'll be using the installed package.
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile && yarn cache clean

FROM node:12-alpine as build

COPY --from=install /node_modules .
COPY . .
RUN yarn build:all:${APP_ENV}


FROM node:12-alpine

# https://github.com/puppeteer/puppeteer/blob/master/docs/troubleshooting.md#running-on-alpine
# Installs latest Chromium package.
RUN apk add --no-cache \
      chromium \
      nss \
      freetype \
      freetype-dev \
      harfbuzz \
      ca-certificates \
      ttf-freefont \
      nodejs \
      yarn

# Tell Puppeteer to skip installing Chrome. We'll be using the installed package.
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

COPY --from=install /node_modules .
COPY . .
RUN yarn build:all:${APP_ENV}



