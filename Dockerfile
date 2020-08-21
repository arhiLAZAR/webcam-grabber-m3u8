FROM ubuntu:bionic
LABEL maintainer="Anton Evdokimov <arhiLAZAR@tandex.ru>"

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
