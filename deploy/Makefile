TERRAFORM?=terraform
TERRAFORM_VERBS=init validate plan apply destroy fmt version

.PHONY: $(TERRAFORM_VERBS)
$(filter-out version init fmt,$(TERRAFORM_VERBS)): version .terraform
$(TERRAFORM_VERBS):
	$(TERRAFORM) $@

.terraform:
	$(MAKE) init

# Do not require interactive confirmation from user
export TF_CLI_ARGS
version init:  TF_CLI_ARGS=
apply destroy: TF_CLI_ARGS=-auto-approve

# Tune Terraform for non-interactive use
export TF_INPUT=0
export TF_IN_AUTOMATION=yes

# Use Russian mirror of Terraform registry
export TF_CLI_CONFIG_FILE=yandex.tfrc

# Yandex Cloud access credentials
ifeq ($(YC_TOKEN),)
MISSING_ENV+=YC_TOKEN
endif
ifeq ($(YC_CLOUD_ID),)
MISSING_ENV+=YC_CLOUD_ID
endif
ifeq ($(YC_FOLDER_ID),)
MISSING_ENV+=YC_FOLDER_ID
endif
ifneq ($(MISSING_ENV),)
$(warning Required variables not defined: $(MISSING_ENV))
endif
