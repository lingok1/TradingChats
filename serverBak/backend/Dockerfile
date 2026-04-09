FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata wget

COPY build/server /app/server
COPY docs /app/docs

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/swagger/index.html >/dev/null || exit 1

CMD ["/app/server"]
