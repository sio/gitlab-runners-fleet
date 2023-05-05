#
# Modify Debian cloud VM image for deploying to Yandex cloud
#

REMOTE?=https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-genericcloud-amd64.qcow2
INPUT?=input.qcow2
SCRATCH?=scratch.raw
OUTPUT?=output.qcow2
COMPACT?=yes

TEMPLATE=$(CURDIR)/template
NSPAWN=systemd-nspawn
NSPAWN+=--image $(SCRATCH)
NSPAWN+=--network-veth
NSPAWN+=--bind="$(NSPAWN_RESULT):/tmp/prepare/result"
NSPAWN+=--bind-ro="$(TEMPLATE):/tmp/prepare/template"
NSPAWN+=--bind-ro="$(TEMPLATE)/prepare.service:/etc/systemd/system/prepare.service"
NSPAWN_RESULT=$(CURDIR)/nspawn.result

LOOP:=$(shell losetup -f)
LOOPp1=/dev/mapper/$(notdir $(LOOP))p1

INTERNET:=$(shell ip route show default | sed -n 's/.* dev \([^\ ]*\) .*/\1/p')
TMP_FIREWALL:=$(shell mktemp -t firewall.XXXXXX)
TMP_FORWARDING:=$(shell mktemp -t forwarding.XXXXXX)

ifeq (yes,$(COMPACT))
QEMU_CONVERT_ARGS+=-c
endif

.PHONY: image
image: $(OUTPUT)

$(INPUT):
	curl -sSL "$(REMOTE)" -o "$@"

$(SCRATCH): $(INPUT)
ifeq (,$(LOOP))
	$(error Required variable not defined: LOOP)
endif
	qemu-img convert -f qcow2 -O raw $(INPUT) $(SCRATCH)
	qemu-img resize -f raw $(SCRATCH) 5G
	fdisk -l $(SCRATCH)
	growpart $(SCRATCH) 1
	losetup $(LOOP) $(SCRATCH)
	kpartx -av $(LOOP)
	e2fsck -fy $(LOOPp1)
	resize2fs $(LOOPp1)
	sync
	kpartx -dv $(LOOP)
	losetup -d $(LOOP)
	fdisk -l $(SCRATCH)

.PHONY: prepare
prepare: $(SCRATCH).done
$(SCRATCH).done: $(SCRATCH)
ifeq (,$(INTERNET))
	$(error Required variable not defined: INTERNET)
endif
	sysctl -n net.ipv4.ip_forward > $(TMP_FORWARDING)
	iptables-save -f $(TMP_FIREWALL)
	systemctl is-active systemd-networkd || systemctl start systemd-networkd
	networkctl
	ip addr
	iptables -t nat -A POSTROUTING -o $(INTERNET) -j MASQUERADE
	iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
	iptables -A FORWARD -i ve-+ -o $(INTERNET) -j ACCEPT
	touch "$(NSPAWN_RESULT)"
	$(NSPAWN) /usr/bin/systemd --unit=prepare.service
	iptables-restore $(TMP_FIREWALL)
	sysctl net.ipv4.ip_forward=$$(cat $(TMP_FORWARDING))
	$(RM) -v $(TMP_FIREWALL) $(TMP_FORWARDING)
	RESULT=$$(cat "$(NSPAWN_RESULT)"); echo "$$RESULT"; [ "$$RESULT" = "success exited 0" ]
	$(RM) -v "$(NSPAWN_RESULT)"
	mv -v $(SCRATCH) $@

$(OUTPUT): $(SCRATCH).done
	qemu-img convert $(QEMU_CONVERT_ARGS) -f raw -O qcow2 -o cluster_size=2M $< $@
	$(RM) -v $<
	ls -lh $@

.PHONY: shell
shell: $(SCRATCH)
	$(NSPAWN) /bin/bash || machinectl shell $(basename $(SCRATCH))

.PHONY: clean
clean:
	-$(RM) -v $(OUTPUT) $(SCRATCH) $(SCRATCH).done
