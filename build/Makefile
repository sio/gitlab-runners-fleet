REMOTE=https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-genericcloud-amd64.qcow2
LOCAL=image.qcow2
DEVICE=/dev/nbd15
MOUNTPOINT=content


$(LOCAL):
	curl -sSL "$(REMOTE)" -o "$@"

$(MOUNTPOINT):
	mkdir -p "$@"

.PHONY: mount
mount: $(LOCAL) $(MOUNTPOINT)
	modprobe nbd max_part=8
	qemu-nbd --connect=$(DEVICE) --format=qcow2 $(LOCAL)
	fdisk -l $(DEVICE)
	mount $(DEVICE)p1 $(MOUNTPOINT)
	ls $(MOUNTPOINT)

.PHONY: umount
umount:
	umount $(MOUNTPOINT)
	qemu-nbd --disconnect $(DEVICE)

.PHONY: rmmod
rmmod:
	rmmod nbd
