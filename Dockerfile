FROM alpine
WORKDIR /www
COPY . .
EXPOSE 1313
CMD ["bin/hugo", "server", "--bind", "0.0.0.0"]