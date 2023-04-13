ARG DEBIAN_VERSION=bullseye-slim
ARG GOLANG_VERSION=1
ARG TERRAFORM_VERSION=1.4


FROM golang:${GOLANG_VERSION} as builder
COPY scale /app
WORKDIR /app
RUN go mod download && go build -o /app/scale


FROM hashicorp/terraform:${TERRAFORM_VERSION} as terraform


FROM debian:${DEBIAN_VERSION}
ENV SCALE_CONFIG=/infra/state/scale.config
ENV SCALE_OUTPUT=/infra/scale.auto.tfvars.json
ENV SCALE_OUTPUT_TMP=/infra/state/scale.auto.tfvars.tmp
COPY --from=builder /app/scale /usr/local/bin/scale
COPY --from=terraform /bin/terraform /usr/local/bin/terraform
RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        ca-certificates \
        make \
    && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
COPY deploy /infra
COPY container/backend.tf /infra/
WORKDIR /infra
RUN echo '{}' > ${SCALE_OUTPUT} && \
    chmod a+w ${SCALE_OUTPUT} && \
    make init && \
    mkdir /infra/state && \
    chmod a+w /infra/state
USER 24601
CMD ["/infra/loop.sh"]
