FROM ubuntu:bionic
LABEL maintainer="Anton Evdokimov <arhiLAZAR@yandex.ru>"

RUN ln -snf /usr/share/zoneinfo/Europe/Moscow /etc/localtime && echo Europe/Moscow > /etc/timezone

RUN apt-get update
RUN apt-get install -y vim tar build-essential devscripts debhelper yasm x264 libx264-dev

RUN \
git clone https://github.com/FFmpeg/FFmpeg && \
cd FFmpeg && \
git checkout n4.3.1 && \
./configure --prefix=/usr/local/ffmpeg --enable-shared --enable-libx264 --enable-gpl && \
make && \
make install && \
echo "export LD_LIBRARY_PATH=/usr/local/ffmpeg/lib" >> ~/.bashrc && \
export LD_LIBRARY_PATH=/usr/local/ffmpeg/lib

ENV LD_LIBRARY_PATH=/usr/local/ffmpeg/lib

RUN cp /FFmpeg/ffmpeg /usr/local/bin/ffmpeg

RUN \
curl -L -O https://golang.org/dl/go1.15.linux-amd64.tar.gz && \
tar -zxf go1.15.linux-amd64.tar.gz && \
rm go*.tar.gz && \
mv go /usr/local/

ENV PATH=$PATH:/usr/local/go/bin

# RUN \
# git clone https://github.com/arhiLAZAR/webcam-grabber-m3u8 && \
# cd webcam-grabber-m3u8 && \
# /usr/local/go/bin/go build . && \
# mv webcam-grabber-m3u8 /root/webcam-grabber-m3u8

COPY webcam-grabber-m3u8 /root/webcam-grabber-m3u8
COPY client_secret.json /root/client_secret.json
COPY youtube-go-quickstart.json /root/.credentials/youtube-go-quickstart.json

WORKDIR /root

ENTRYPOINT /root/webcam-grabber-m3u8
