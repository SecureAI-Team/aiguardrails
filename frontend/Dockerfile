FROM node:20-alpine AS build
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci
COPY . .
ARG VITE_API_BASE=http://api:8080
ARG VITE_ADMIN_TOKEN=changeme-admin
ENV VITE_API_BASE=${VITE_API_BASE}
ENV VITE_ADMIN_TOKEN=${VITE_ADMIN_TOKEN}
RUN npm run build

FROM nginx:1.25-alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

