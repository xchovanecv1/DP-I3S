FROM ubuntu:16.04

# Install build dependencies (and vim + picocom for editing/debugging)
RUN apt-get -qq update \
    && apt-get install -y gcc git wget make libncurses-dev flex bison gperf python python-serial python-pip \
                          cmake ninja-build \
                          ccache \
                          vim picocom \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Get the ESP32 toolchain
ENV ESP_TCHAIN_BASEDIR /opt/local/espressif

RUN mkdir -p $ESP_TCHAIN_BASEDIR \
    && wget -O $ESP_TCHAIN_BASEDIR/esp32-toolchain.tar.gz \
            https://dl.espressif.com/dl/xtensa-esp32-elf-linux64-1.22.0-80-g6c4433a-5.2.0.tar.gz \
    && tar -xzf $ESP_TCHAIN_BASEDIR/esp32-toolchain.tar.gz \
           -C $ESP_TCHAIN_BASEDIR/ \
    && rm $ESP_TCHAIN_BASEDIR/esp32-toolchain.tar.gz

RUN mkdir -p $ESP_TCHAIN_BASEDIR \
    && wget -O $ESP_TCHAIN_BASEDIR/esp32ulp-toolchain.tar.gz \
            https://dl.espressif.com/dl/esp32ulp-elf-binutils-linux64-d2ae637d.tar.gz \
    && tar -xzf $ESP_TCHAIN_BASEDIR/esp32ulp-toolchain.tar.gz \
           -C $ESP_TCHAIN_BASEDIR/ \
    && rm $ESP_TCHAIN_BASEDIR/esp32ulp-toolchain.tar.gz

# Setup IDF_PATH
ENV IDF_PATH /esp/esp-idf
ENV MDF_PATH /esp/esp-mdf

COPY ./esp-mdf/esp-idf /esp/esp-idf
COPY ./esp-mdf /esp/esp-mdf

RUN /usr/bin/python -m pip install --user -r /esp/esp-mdf/esp-idf/requirements.txt

#RUN mkdir -p $IDF_PATH
#RUN mkdir -p $MDF_PATH

# Add the toolchain binaries to PATH
ENV PATH $ESP_TCHAIN_BASEDIR/xtensa-esp32-elf/bin:$ESP_TCHAIN_BASEDIR/esp32ulp-elf-binutils/bin:$IDF_PATH/tools:$PATH

# This is the directory where our project will show up
RUN mkdir -p /esp/project
WORKDIR /esp/project
ENTRYPOINT ["/bin/bash"]
