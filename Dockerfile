FROM node:18-alpine as web-builder
WORKDIR /web
COPY web/package.json web/pnpm-lock.yaml ./
COPY web/ ./
RUN npm install -g pnpm && pnpm install
RUN npm run build

FROM golang:1.22 as server-builder
WORKDIR /server
COPY . .
RUN CGO_ENABLED=0 go build -o chzzk-live-dl .

FROM alpine:latest
WORKDIR /root/
COPY --from=server-builder /server/chzzk-live-dl .
COPY --from=web-builder /web/dist ./public
EXPOSE 2000

CMD ["./chzzk-live-dl", "--dir=/streams"]

