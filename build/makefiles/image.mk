#
# Modify Debian cloud VM image for deploying to Yandex cloud
#

REMOTE=https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-genericcloud-amd64.qcow2
INPUT=input.qcow2
OUTPUT=output.qcow2
DEVICE=/dev/nbd15
MOUNTPOINT=/tmp/gitlab-runners-fleet-rootfs
TEMPLATE=template
SCRIPT=install.sh

.PHONY: image
image: mount chroot umount

$(INPUT):
	curl -sSL "$(REMOTE)" -o "$@"

$(OUTPUT): $(INPUT)
	cp $(INPUT) $(OUTPUT)
	qemu-img resize $(OUTPUT) 5G


$(MOUNTPOINT):
	mkdir -p "$@"

.PHONY: clean
clean:
	$(RM) -v $(OUTPUT)

.PHONY: mount
mount: $(OUTPUT) $(MOUNTPOINT)
	modprobe nbd max_part=8
	qemu-nbd --connect=$(DEVICE) --format=qcow2 $(OUTPUT)
	-growpart $(DEVICE) 1
	e2fsck -fy $(DEVICE)p1
	resize2fs $(DEVICE)p1
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
	mkdir -p $(MOUNTPOINT)/etc/provision
	cp -avr $(TEMPLATE)/* $(MOUNTPOINT)/etc/provision
	chroot $(MOUNTPOINT) /etc/provision/$(SCRIPT)

.PHONY: compact
compact:
	qemu-img convert -c -f qcow2 -O qcow2 -o cluster_size=2M $(OUTPUT) $(OUTPUT).compact
	ls -lh $(INPUT) $(OUTPUT)*
	mv $(OUTPUT).compact $(OUTPUT)
