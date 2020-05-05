FROM adoptopenjdk/openjdk10:alpine-slim

ENV LANG=C.UTF-8 \
    JAVA_HOME=/opt/java/openjdk \
    PATH=${PATH}:/opt/java/openjdk/bin \
    LANG=C.UTF-8 \
    KAFKA_VERSION="2.4.0"

RUN sed -i s/#networkaddress.cache.ttl=-1/networkaddress.cache.ttl=10/ $JAVA_HOME/conf/security/java.security

RUN mkdir -p /opt/kafka
WORKDIR /opt/kafka

RUN \
 mkdir -p /opt/kafka && \
 cd /opt/kafka && \
 apk add --no-cache wget ca-certificates bash && \
 wget http://apache.mirror.anlx.net/kafka/${KAFKA_VERSION}/kafka_2.12-${KAFKA_VERSION}.tgz && \
 md5sum kafka_2.12-${KAFKA_VERSION}.tgz && \
 tar -zxvf kafka_2.12-${KAFKA_VERSION}.tgz --strip 1 && \
 rm kafka_2.12-${KAFKA_VERSION}.tgz && \
 apk del wget

RUN ["/bin/bash", "-c", "echo \"export LANG=C.UTF-8\" > /etc/profile.d/locale.sh"]
WORKDIR /opt/kafka/bin

CMD ["./kafka-server-start.sh", "../config/server.properties"]
