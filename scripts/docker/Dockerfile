FROM alpine
RUN apk update \
    apk add busybox-extras \
    apk add curl

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
ARG APPID
ENV APPID ${APPID}
# set different CUSTOM_RUNTIME_ENV using config/xxx.yaml
ENV CUSTOM_RUNTIME_ENV "pro"
ENV IN_DOCKER "true"
ADD ./$APPID /root/$APPID
WORKDIR /root/$APPID
RUN chmod a+x /root/$APPID/entrypoint.sh
EXPOSE 8080
CMD /root/$APPID/entrypoint.sh /root/$APPID