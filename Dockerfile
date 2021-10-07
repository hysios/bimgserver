FROM golang:1.17.1

WORKDIR /go/src/app

COPY . .

COPY ./etc/sources.list /etc/apt/sources.list

RUN wget http://keyue-cloud.oss-cn-guangzhou.aliyuncs.com/downloads%2Fvips-8.11.4.tar.gz -O /tmp/vips-8.11.4.tar.gz

RUN apt-get update -y

RUN apt-get --fix-missing install glib2.0 libexpat1-dev libjpeg-dev libpng-dev libexif-dev -y

RUN cd /tmp; tar xf vips-8.11.4.tar.gz; cd vips-8.11.4; ./configure

WORKDIR /tmp/vips-8.11.4

RUN make; make install; ldconfig

ENV GOPROXY=https://goproxy.cn

WORKDIR /go/src/app

CMD ["go", "build", "."]