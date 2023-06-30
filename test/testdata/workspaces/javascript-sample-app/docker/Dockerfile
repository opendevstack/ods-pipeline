FROM node:16-alpine

COPY src /node/
COPY node_modules /node/
COPY package.json /node/
COPY package-lock.json /node/

RUN chown -R node: /node

EXPOSE 8080

USER node

WORKDIR /node

CMD ["node", "src/index.js"]
