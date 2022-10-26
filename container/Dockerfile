FROM python:3.10-slim

ARG PULUMI_VERSION=3.29.1
ARG PULUMI_YANDEX_VERSION=0.13.0
ARG PULUMI_PLATFORM=linux-x64
ENV PULUMI_HOME=/opt/pulumi/x/

ADD https://get.pulumi.com/releases/sdk/pulumi-v${PULUMI_VERSION}-${PULUMI_PLATFORM}.tar.gz /pulumi.tar.gz
ADD . /src
ADD container/entrypoint.sh /entrypoint

RUN apt-get update && apt-get install -y libnss-wrapper && \
    apt-get clean && rm -rf /var/lib/apt/lists/* && \
    mkdir -p /opt && \
    tar -xzvf /pulumi.tar.gz -C /opt && \
    rm -vf \
        /opt/pulumi/pulumi-language-dotnet \
        /opt/pulumi/pulumi-language-go \
        /opt/pulumi/pulumi-language-java \
        /opt/pulumi/pulumi-language-nodejs \
        /opt/pulumi/pulumi-language-yaml \
        /pulumi.tar.gz && \
    /opt/pulumi/pulumi plugin install \
        --non-interactive \
        --exact \
        resource yandex ${PULUMI_YANDEX_VERSION} && \
    chmod a+rwX "${PULUMI_HOME}" "${PULUMI_HOME}/plugins" && \
    chmod a+rX -R "${PULUMI_HOME}/plugins" && \
    useradd --no-create-home --shell /usr/sbin/nologin fleetmanager && \
    mkdir -p /home/fleetmanager && \
    chown fleetmanager /home/fleetmanager && \
    SETUPTOOLS_SCM_PRETEND_VERSION=999 \
    pip --no-cache-dir install /src && \
    rm -rf /src && \
    mkdir -p /state && \
    chmod 777 /state

USER fleetmanager

WORKDIR /config
ENTRYPOINT /entrypoint
