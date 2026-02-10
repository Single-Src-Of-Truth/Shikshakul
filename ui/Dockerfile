FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --legacy-peer-deps

COPY . .

# 1. Admin -> /admin/
RUN npx nx run administrator-portal:build:production --base-href /admin/

# 2. Teacher -> /teacher/
RUN npx nx run teacher-portal:build:production --base-href /teacher/

# 3. Guardian -> /parent/ (URL is parent, app name is guardian-portal)
RUN npx nx run guardian-portal:build:production --base-href /parent/


FROM nginx:alpine

RUN rm -rf /usr/share/nginx/html/*

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=builder /app/dist/apps/administrator-portal/browser /usr/share/nginx/html/admin
COPY --from=builder /app/dist/apps/teacher-portal/browser       /usr/share/nginx/html/teacher
COPY --from=builder /app/dist/apps/guardian-portal/browser      /usr/share/nginx/html/guardian

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
