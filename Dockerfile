FROM alpine as build
WORKDIR /app/
ADD https://github.com/caddyserver/caddy/releases/download/v2.3.0/caddy_2.3.0_linux_amd64.tar.gz caddy.tar.gz
RUN tar -xvf caddy.tar.gz
RUN apk add hugo
COPY . .
RUN hugo

FROM busybox as production
COPY --from=build /app/caddy /bin
COPY --from=build /app/public /var/www/html
RUN adduser -D caddy
RUN adduser caddy caddy
USER caddy
WORKDIR /home/caddy
ENTRYPOINT ["/bin/caddy", "file-server", "-root", "/var/www/html", "-browse", "-listen", ":8080", "-access-log"]
