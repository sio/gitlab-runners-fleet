REMOTE=https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-genericcloud-amd64.qcow2
INPUT=input.qcow2
OUTPUT=output.qcow2
DEVICE=/dev/nbd15
MOUNTPOINT=rootfs
TEMPLATE=template
SCRIPT=install.sh

.PHONY: image
image: mount chroot umount

$(INPUT):
	curl -sSL "$(REMOTE)" -o "$@"

$(OUTPUT): $(INPUT)
	cp $(INPUT) $(OUTPUT)

$(MOUNTPOINT):
	mkdir -p "$@"

.PHONY: mount
mount: $(OUTPUT) $(MOUNTPOINT)
	modprobe nbd max_part=8
	qemu-nbd --connect=$(DEVICE) --format=qcow2 $(OUTPUT)
	fdisk -l $(DEVICE)
	mount $(DEVICE)p1 $(MOUNTPOINT)
	mount --bind /dev $(MOUNTPOINT)/dev
	mount --bind /dev/pts $(MOUNTPOINT)/dev/pts
	mount --bind /proc $(MOUNTPOINT)/proc
	mount --bind /sys $(MOUNTPOINT)/sys
	mount --bind /run $(MOUNTPOINT)/run
	mv -v $(MOUNTPOINT)/etc/resolv.conf $(MOUNTPOINT)/etc/resolv.conf.orig
	cat /etc/resolv.conf > $(MOUNTPOINT)/etc/resolv.conf
	ls $(MOUNTPOINT)
	tail -n+0 -v $(MOUNTPOINT)/etc/*release*

.PHONY: umount
umount:
	mv -v $(MOUNTPOINT)/etc/resolv.conf.orig $(MOUNTPOINT)/etc/resolv.conf
	-umount $(MOUNTPOINT)/proc
	-umount $(MOUNTPOINT)/dev/pts
	-umount $(MOUNTPOINT)/dev
	-umount $(MOUNTPOINT)/sys
	-umount $(MOUNTPOINT)/run
	umount $(MOUNTPOINT)
	qemu-nbd --disconnect $(DEVICE)
	-rmmod nbd
	rmdir $(MOUNTPOINT)

.PHONY: chroot
chroot:
	cp -avr $(TEMPLATE) $(MOUNTPOINT)/etc/provision
	chroot $(MOUNTPOINT) /etc/provision/$(SCRIPT)

.PHONY: compact
compact:
	mv $(OUTPUT) $(OUTPUT).orig
	qemu-img convert -c -f qcow2 -O qcow2 $(OUTPUT).orig $(OUTPUT)
	ls -l $(OUTPUT)*
	rm $(OUTPUT).orig

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

.PHONY: .gha-environment
.gha-environment:
	apt update
	apt-get install -y qemu-utils

include ../deploy/Makefile
.PHONY: bucket
bucket: apply
	$(RM) $(BUCKET_EXTRAS)
export TF_VAR_ycs3_bucket=$(S3_BUCKET)
BUCKET_EXTRAS=yandex.tf yandex.tfrc
$(TERRAFORM_VERBS): $(BUCKET_EXTRAS)
$(BUCKET_EXTRAS):
	cp ../deploy/$@ $@
