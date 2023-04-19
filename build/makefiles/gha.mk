#
# Prepare GitHub actions environment for building VM image
#

.PHONY: .gha-environment
.gha-environment:
	apt update
	apt-get install -y qemu-utils kpartx systemd-container
