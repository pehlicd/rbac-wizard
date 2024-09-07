FROM node:22 AS frontend-builder

WORKDIR /app/ui

COPY ui /app/ui

RUN npm install -g npm@latest && \
    npm install --force
RUN npm run build

FROM golang:1.23 as backend-builder

WORKDIR /app

COPY . .

RUN go install github.com/rakyll/statik@latest

COPY --from=frontend-builder /app/ui/dist /app/ui/dist

RUN statik -src=./ui/dist/ -dest=./internal/ -f
RUN go build -o rbac-wizard

FROM alpine:3.20.2

COPY --from=backend-builder /app/rbac-wizard /usr/local/bin/rbac-wizard

RUN apk add libc6-compat

EXPOSE 8080

ENTRYPOINT ["rbac-wizard"]

CMD ["serve"]