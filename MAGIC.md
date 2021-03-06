# TODO: get rid of magic values used in the code

Currently some values are hardcoded in multiple places without an easy way to
change them. This is fine for a target audience of a single person (me), but
may not be desirable for other users.

These values should be refactored into constants with a single source of truth
to allow changing:

## High priority magic values

- `private-runner` - GitLab runner tag
- `9252` - GitLab runner metrics port on localhost
- `8080` - GitLab runner metrics port on WAN

## Low priority magic values

- `https://gitlab.com/sio/...` URLs
- `https://github.com/sio/...` URLs
