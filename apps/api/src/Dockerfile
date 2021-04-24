FROM node:14

WORKDIR /app

COPY package.json ./
RUN npm install --only=production

COPY . ./
CMD npm run start