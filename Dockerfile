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

RUN \
curl -L -O https://golang.org/dl/go1.15.linux-amd64.tar.gz && \
tar -zxf go1.15.linux-amd64.tar.gz && \
rm go*.tar.gz && \
mv go /usr/local/ && \
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && \
export PATH=$PATH:/usr/local/go/bin

RUN \
git clone https://github.com/arhiLAZAR/webcam-grabber-m3u8 && \
cd webcam-grabber-m3u8 && \
go build . && \
mv webcam-grabber-m3u8 /usr/local/bin/webcam-grabber-m3u8 && \
mkdir /etc/webcam-grabber-m3u8 && \
cp debian/webcam-grabber-m3u8.service /lib/systemd/system/webcam-grabber-m3u8.service && \
systemctl daemon-reload && \
systemctl enable webcam-grabber-m3u8.service

# && systemctl start webcam-grabber-m3u8.service
