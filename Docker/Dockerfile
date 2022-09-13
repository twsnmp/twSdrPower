FROM alpine:latest as builder
RUN apk add --update --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata
COPY twSdrPower /
RUN  chmod +x /twSdrPower

FROM alpine:latest
COPY --from=builder /twSdrPower /
COPY --from=builder /etc/timezone /etc/
COPY --from=builder /etc/localtime /etc/
WORKDIR /
ENTRYPOINT [ "/twSdrPower" ]

