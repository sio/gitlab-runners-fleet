# Install Pulumi and required packages on management node (Debian)


ifdef PULUMI_AUTO_INSTALL

ifneq (ok,$(shell $(PULUMI) version >/dev/null 2>&1 && echo ok))
PULUMI=$(PULUMI_PROJECT)/bin/pulumi
PULUMI_DOWNLOAD_VERSION?=2.22.0
endif

check-software: management-node-packages $(PULUMI)

$(PULUMI_PROJECT)/bin/pulumi:
	mkdir -p "$(dir $@)"
	curl --fail -sSL "https://get.pulumi.com/releases/sdk/pulumi-v$(PULUMI_DOWNLOAD_VERSION)-linux-x64.tar.gz" \
	| tar zxv -C "$(dir $@)" --strip-components=1

.PHONY: management-node-packages
management-node-packages: $(PULUMI_PROJECT)/bin/os_packages_are_installed

$(PULUMI_PROJECT)/bin/os_packages_are_installed:
	apt-get -y install make python3-venv openssh-client curl
	mkdir -p $(dir $@)
	touch $@

endif # PULUMI_AUTO_INSTALL
