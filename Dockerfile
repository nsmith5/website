FROM alpine
WORKDIR /website
COPY public .
COPY bin/hugo .
EXPOSE 1313
CMD ["hugo", "serve"]