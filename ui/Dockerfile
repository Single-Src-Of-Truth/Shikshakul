FROM oven/bun:alpine AS builder

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .

# Admin -> /admin/
RUN bunx nx run administrator-portal:build:production --base-href /admin/

# Teacher -> /teacher/
RUN bunx nx run teacher-portal:build:production --base-href /teacher/

# Guardian -> /parent/
RUN bunx nx run guardian-portal:build:production --base-href /parent/


FROM nginx:alpine

RUN rm -rf /usr/share/nginx/html/*

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=builder /app/dist/apps/administrator-portal/browser /usr/share/nginx/html/admin
COPY --from=builder /app/dist/apps/teacher-portal/browser       /usr/share/nginx/html/teacher
COPY --from=builder /app/dist/apps/guardian-portal/browser      /usr/share/nginx/html/parent

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
