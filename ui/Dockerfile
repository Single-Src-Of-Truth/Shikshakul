FROM oven/bun:alpine AS builder

WORKDIR /app

ENV NODE_ENV=production

COPY package.json bun.lock ./
RUN --mount=type=cache,id=bun,target=/root/.bun-cache \
    bun install --frozen-lockfile

COPY . .

# Build all portals in parallel using Nx and cache mounts for speed
RUN --mount=type=cache,target=/app/.nx/cache \
    bunx nx run administrator-portal:build:production --base-href /admin/ && \
    bunx nx run teacher-portal:build:production --base-href /teacher/ && \
    bunx nx run guardian-portal:build:production --base-href /parent/


FROM nginx:alpine

RUN rm -rf /usr/share/nginx/html/*

COPY nginx.conf /etc/nginx/conf.d/default.conf

COPY --from=builder /app/dist/apps/administrator-portal/browser /usr/share/nginx/html/admin
COPY --from=builder /app/dist/apps/teacher-portal/browser       /usr/share/nginx/html/teacher
COPY --from=builder /app/dist/apps/guardian-portal/browser      /usr/share/nginx/html/parent

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
