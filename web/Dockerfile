FROM node:16.13-alpine as builder

WORKDIR /app

COPY . .
RUN npm install
RUN npm run build

FROM node:16.13-alpine
WORKDIR /app
COPY package.json .
COPY --from=builder /app/build build
RUN npm install --only=prod

CMD ["node", "build"]
