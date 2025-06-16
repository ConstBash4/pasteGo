FROM node:alpine AS frontend-builder
WORKDIR /frontend
#ADD ./frontend/package.json .
#ADD ./frontend/package-lock.json .
COPY ./frontend .
RUN npm ci && npm run build
#RUN npm run build

FROM golang:alpine AS backend-builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN apk add --no-cache build-base libc-dev
RUN CGO_ENABLED=1 GOOS=linux go build -o main main.go 

FROM alpine
WORKDIR /app/build
COPY --from=frontend-builder /frontend/build /app/build/
WORKDIR /app
COPY --from=backend-builder /build/main /app/main
COPY .env /app/.env

EXPOSE 10015

CMD ["./main"]