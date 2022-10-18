# Install Pulumi and required packages on management node (Debian)


ifdef PULUMI_INSTALL_AUTO

ifneq (ok,$(shell $(PULUMI) version >/dev/null 2>&1 && echo ok))
PULUMI_INSTALL_VERSION?=3.29.1
PULUMI_INSTALL_ARCH?=linux-x64
PULUMI_INSTALL_DIR?=./bin
PULUMI=$(PULUMI_INSTALL_DIR)/pulumi
endif

check-software: $(PULUMI)

$(PULUMI):
	mkdir -p "$(dir $@)"
	curl --fail -sSL "https://get.pulumi.com/releases/sdk/pulumi-v$(PULUMI_INSTALL_VERSION)-$(PULUMI_INSTALL_ARCH).tar.gz" \
	| tar zxv -C "$(dir $@)" --strip-components=1

.PHONY: debian-packages
debian-packages: $(dir $(PULUMI))/.debian-packages-are-installed
$(dir $(PULUMI))/.debian-packages-are-installed:
	apt-get -y install make python3-venv openssh-client curl
	mkdir -p $(dir $@)
	touch $@


endif # PULUMI_INSTALL_AUTO
