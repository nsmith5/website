FROM alpine as web
WORKDIR /workspace
RUN apk add --no-cache hugo go
COPY . .
RUN hugo

FROM golang as server
WORKDIR /workspace
COPY src ./src
RUN CGO_ENABLED=0 go build -o server src/main.go

FROM scratch
COPY --from=web /workspace/public /public
COPY --from=server /workspace/server /
EXPOSE 3000
CMD ["/server", "-dir", "public"]
