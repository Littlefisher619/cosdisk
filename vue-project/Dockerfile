FROM node:16-alpine

# ENV GQL_BACKEND 
ARG GQL_BACKEND=http://localhost:8080/query

RUN echo build with GQL_BACKEND=$GQL_BACKEND
WORKDIR /app

COPY . .
RUN npm config set registry http://mirrors.cloud.tencent.com/npm/
RUN npm install
RUN echo "VUE_APP_GQL_BACKEND=$GQL_BACKEND" >> .env
RUN npm run build

FROM nginx:1.17
WORKDIR /usr/share/nginx/html
COPY --from=0 /app/dist .