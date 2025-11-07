FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./frontend/
RUN --mount=type=cache,target=/root/.npm cd frontend && npm ci --prefer-offline --no-audit --silent
COPY frontend ./frontend
RUN cd frontend && npm run build

FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod ./backend/
RUN --mount=type=cache,target=/go/pkg/mod cd backend && go mod download
COPY backend ./backend
RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o rater .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=backend-builder /app/backend/rater /rater
COPY --from=frontend-builder /app/frontend/dist /frontend/dist
ENV PORT=8080 STATIC_DIR=/frontend/dist
EXPOSE 8080
ENTRYPOINT ["/rater"]

