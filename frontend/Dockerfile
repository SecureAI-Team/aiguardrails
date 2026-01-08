FROM node:20-alpine AS build
WORKDIR /app
# 使用淘宝npm镜像加速
RUN npm config set registry https://registry.npmmirror.com
COPY package.json ./
RUN npm install
COPY . .
ARG VITE_API_BASE=/v1
ARG VITE_ADMIN_TOKEN=changeme-admin
ENV VITE_API_BASE=${VITE_API_BASE}
ENV VITE_ADMIN_TOKEN=${VITE_ADMIN_TOKEN}
RUN npm run build

FROM nginx:1.25-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]



