# Must be provided by user
REQUIRED_ENVIRONMENT_VARIABLES= \
	GITLAB_API_TOKEN \
	GITLAB_RUNNER_TOKEN \
	HCLOUD_TOKEN \
	PULUMI_CONFIG_PASSPHRASE \

# Sane defaults are provided for these
REQUIRED_ENVIRONMENT_VARIABLES+= \
	GITLAB_RUNNER_SSHKEY \
	HCLOUD_USERDATA_TEMPLATE \
	PULUMI_SNAPSHOT_OBJECT \

define require-env
ifndef $(1)
MISSING_ENVIRONMENT_VARIABLES+=$(1)
endif
endef

$(foreach env,$(REQUIRED_ENVIRONMENT_VARIABLES),$(eval $(call require-env, $(env))))

ifneq (,$(strip $(MISSING_ENVIRONMENT_VARIABLES)))
$(error Required environment variables not defined: $(MISSING_ENVIRONMENT_VARIABLES))
endif
