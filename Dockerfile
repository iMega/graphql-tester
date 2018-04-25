FROM scratch

MAINTAINER Dmitry Stoletoff <info@imega.ru>

WORKDIR /data
ADD build/rootfs.tar.gz /

ENTRYPOINT ["/usr/bin/graphql-tester"]
