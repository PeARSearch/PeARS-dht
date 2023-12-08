FROM gitpod/workspace-full:latest

RUN sudo apt -y update

RUN sudo sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

# Install OpenDHT dependencies
RUN sudo apt-get install -y \
        dialog apt-utils \
    && sudo apt-get clean \
    && sudo echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections

RUN sudo apt-get update && sudo apt-get install -y \
        build-essential pkg-config cmake git wget \
        libtool autotools-dev autoconf graphviz doxygen\
        cython3 python3-dev python3-setuptools python3-build python3-virtualenv \
        libncurses5-dev libreadline-dev nettle-dev libcppunit-dev \
        libgnutls28-dev libuv1-dev libjsoncpp-dev libargon2-dev \
        libssl-dev libfmt-dev libhttp-parser-dev libasio-dev libmsgpack-dev  openssh-client \
    && sudo apt-get clean && sudo  rm -rf /var/lib/apt/lists/* /var/cache/apt/*

# Build & install restinio (for proxy server/client):
RUN mkdir restinio && cd restinio \
    && wget https://github.com/aberaud/restinio/archive/2224ffedef52cb2b74645d63d871d61dbd0f165e.tar.gz \
    && ls -l && tar -xzf 2224ffedef52cb2b74645d63d871d61dbd0f165e.tar.gz \
    && cd restinio-2224ffedef52cb2b74645d63d871d61dbd0f165e/dev \
    && sudo cmake -DCMAKE_INSTALL_PREFIX=/usr -DRESTINIO_TEST=OFF -DRESTINIO_SAMPLE=OFF \
             -DRESTINIO_INSTALL_SAMPLES=OFF -DRESTINIO_BENCH=OFF -DRESTINIO_INSTALL_BENCHES=OFF \
             -DRESTINIO_FIND_DEPS=ON -DRESTINIO_ALLOW_SOBJECTIZER=Off -DRESTINIO_USE_BOOST_ASIO=none . \
    && make -j8 && sudo make install \
    && cd ../../../ && sudo rm -rf restinio

RUN git clone https://github.com/nandajavarma/opendht.git /workspace/opendht

RUN sudo chown -hR gitpod:gitpod /workspace

WORKDIR /workspace/opendht

RUN ls -a
# build and install
RUN cmake -DCMAKE_INSTALL_PREFIX=/usr \
				-DCMAKE_INTERPROCEDURAL_OPTIMIZATION=On \
				-DOPENDHT_C=On \
				-DOPENDHT_PEER_DISCOVERY=On \
				-DOPENDHT_PYTHON=On \
				-DOPENDHT_TOOLS=On \
				-DOPENDHT_PROXY_SERVER=On \
				-DOPENDHT_PROXY_CLIENT=On \
            -DOPENDHT_SYSTEMD=Off

RUN pip3 install --upgrade cython
RUN make -j8
RUN sudo make install

ENV PYTHONPATH=/workspace/PeARS-dht

WORKDIR /workspace/PeARS-dht
