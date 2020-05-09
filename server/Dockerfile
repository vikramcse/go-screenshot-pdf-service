FROM golang:1.14-alpine as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o capture-screen .

FROM alpine:3.11 
RUN apk add --no-cache \
      chromium \
      nss \
      freetype \
      freetype-dev \
      harfbuzz \
      ca-certificates \
      ttf-freefont
 
WORKDIR /root/
COPY --from=builder /app/capture-screen .
EXPOSE 8001
CMD ["./capture-screen", "-host", "0.0.0.0", "-port", "8001"]


