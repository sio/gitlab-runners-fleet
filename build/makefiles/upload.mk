.PHONY: upload
AWS_DEFAULT_REGION?=ru-central1
AWS_ENDPOINT?=https://storage.yandexcloud.net
S3=aws s3 --endpoint-url=$(AWS_ENDPOINT)
S3_OBFUSCATION_PREFIX?=$(firstword $(shell echo $(OUTPUT)|md5sum))
export S3_OBFUSCATION_PREFIX
export S3_BUCKET
upload: $(OUTPUT)
ifeq (,$(S3_BUCKET))
	$(error Variable not defined: S3_BUCKET)
endif
	$(S3) cp $(OUTPUT) s3://$$S3_BUCKET/$$S3_OBFUSCATION_PREFIX/base.qcow2
