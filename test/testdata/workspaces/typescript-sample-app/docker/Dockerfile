FROM node:16-alpine

# Might be better to have the npm install done during docker build.
# - node_modules from the build in this file instead of getting the one from the build, which also has test modules.

COPY node_modules /node/node_modules/
COPY package.json /node/
COPY package-lock.json /node/
COPY dist /node/

RUN chown -R node: /node

EXPOSE 8080

USER node

WORKDIR /node

CMD ["node", "src/index.js"]
