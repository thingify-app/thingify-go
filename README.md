# thingify-go

This is a Go implementation of `thingify`, intended to perform the "thing" role (i.e. running onboard a remote-controlled device).

## Initial setup
As this repo depends on private GitHub repos, the following initial setup is required:

- Assuming SSH keys have been setup with access to the `github.com/thingify-app/` organisation, we want Go (or any other Git user) to default to SSH instead of HTTPS for `github.com` URLs. Modify `~/.gitconfig` to add the following lines:
  ```
  [url "ssh://git@github.com/"]
  	insteadOf = https://github.com/
  ```
- Configure Go to treat all repos within `github.com/thingify-app/` as private. This just disables the use of the central Go package proxy:
  ```
  go env -w GOPRIVATE=github.com/thingify-app/*
  ```
