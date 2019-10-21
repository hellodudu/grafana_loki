FROM alpine
RUN apk add tzdata

RUN mkdir /app
WORKDIR /app

COPY loki_conn /app/loki_conn

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' >/etc/timezone 

ENTRYPOINT ["/app/loki_conn"]
