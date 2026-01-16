FROM node:24-alpine

RUN npm install -g pnpm

WORKDIR /opt/app

COPY package.json .

RUN pnpm install

COPY . .

EXPOSE 8080

CMD [ "pnpm", "run", "dev" ]