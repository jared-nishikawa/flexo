FROM golang:latest

RUN apt-get update
RUN apt-get install -y mingw-w64
RUN apt-get install -y cmake patch libxml2-dev libssl-dev xz-utils clang libz-dev
RUN \
    cd / && \
    git clone https://github.com/tpoechtrager/osxcross && \
    cd osxcross && \
    wget https://s3.dockerproject.org/darwin/v2/MacOSX10.11.sdk.tar.xz && \
    mv MacOSX10.11.sdk.tar.xz tarballs && \
    sed -i -e 's|-march=native||g' build_clang.sh wrapper/build_wrapper.sh && \
    #mkdir -p /usr/local/osx-ndk-x86 && \
    UNATTENDED=yes OSX_VERSION_MIN=10.7 ./build.sh && \
    cp -r target/* /usr

RUN dpkg --add-architecture i386 && apt-get update
RUN apt-get install -y xorg-dev libgl1-mesa-dev libgl1-mesa-dev:i386 libxrandr-dev:i386 libxinerama-dev:i386 libxi-dev:i386 libxcursor-dev:i386 libgl1-mesa-dev:i386
Run apt-get install -y gcc-multilib
