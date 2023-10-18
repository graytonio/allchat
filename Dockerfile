FROM alpine:3.18

WORKDIR /app

COPY allchat /app
COPY templates /app/templates
COPY assets /app/assets

RUN touch /config.yaml

ENTRYPOINT [ "/app/allchat" ]