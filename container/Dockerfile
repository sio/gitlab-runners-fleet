FROM python:slim

ARG PULUMI_VERSION=3.29.1
ARG PULUMI_YANDEX_VERSION=0.13.0
ARG PULUMI_PLATFORM=linux-x64

ADD https://get.pulumi.com/releases/sdk/pulumi-v${PULUMI_VERSION}-${PULUMI_PLATFORM}.tar.gz /pulumi.tar.gz
RUN mkdir -p /opt && \
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
        resource yandex ${PULUMI_YANDEX_VERSION}

ADD . /src
RUN SETUPTOOLS_SCM_PRETEND_VERSION=999 \
    pip --no-cache-dir install /src && \
    rm -rf /src && \
    mkdir -p /state && \
    chmod 777 /state

WORKDIR /config
COPY container/entrypoint.sh /entrypoint
ENTRYPOINT /entrypoint
