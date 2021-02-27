REQUIRED_ENVIRONMENT_VARIABLES= \
	GITLAB_API_TOKEN \
	HCLOUD_TOKEN \
	PULUMI_CONFIG_PASSPHRASE

define require-env
ifndef $(1)
MISSING_ENVIRONMENT_VARIABLES+=$(1)
endif
endef

$(foreach env,$(REQUIRED_ENVIRONMENT_VARIABLES),$(eval $(call require-env, $(env))))

ifneq (,$(strip $(MISSING_ENVIRONMENT_VARIABLES)))
$(error Required environment variables not defined: $(MISSING_ENVIRONMENT_VARIABLES))
endif
