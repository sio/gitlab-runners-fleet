.PHONY: upload
AWS_ENDPOINT?=https://storage.yandexcloud.net
export AWS_DEFAULT_REGION?=ru-central1
export AWS_EC2_METADATA_DISABLED=true
S3=aws s3 --endpoint-url=$(AWS_ENDPOINT)
S3_OBFUSCATION_PREFIX?=$(firstword $(shell echo $(OUTPUT)|md5sum))
S3_PATH?=s3://$$S3_BUCKET/$$S3_OBFUSCATION_PREFIX/base.qcow2
export S3_OBFUSCATION_PREFIX
export S3_BUCKET
upload:
ifeq (,$(S3_BUCKET))
	$(error Variable not defined: S3_BUCKET)
endif
	$(S3) cp $(OUTPUT) "$(S3_PATH)" --only-show-errors


.PHONY: remove
destroy: remove
remove:
	-$(S3) rm "$(S3_PATH)"
