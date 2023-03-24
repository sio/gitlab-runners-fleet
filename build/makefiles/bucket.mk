#
# Create object storage bucket for storing VM images
#

include ../deploy/Makefile
.PHONY: bucket
bucket: apply
	$(RM) $(BUCKET_EXTRAS)
export TF_VAR_ycs3_bucket=$(S3_BUCKET)
BUCKET_EXTRAS=yandex.tf yandex.tfrc
.INTERMEDIATE: $(BUCKET_EXTRAS)
$(TERRAFORM_VERBS): $(BUCKET_EXTRAS)
$(BUCKET_EXTRAS):
	cp ../deploy/$@ $@
