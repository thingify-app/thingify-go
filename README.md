# thingify-go

This is a Go implementation of `thingify`, intended to perform the "thing" role (i.e. running onboard a remote-controlled device).

## Docker build
To use the Docker build, which installs the necessary Raspberry Pi dependencies
and cross-compiles for armhf, first build the Docker image:
```
cd build
docker build -t thingify-go-build .
```

Thereafter, the Docker image can be used to build the program. We mount the
`~/.ssh` directory so that the container has access to the necessary private
GitHub repos (note that this assumes the SHH keys are *not* password-protected).
We also mount the `GOCACHE` and `GOMODCACHE` to speed up builds:
```
docker run --rm \
    --mount type=bind,source="$(pwd)",target=/build \
    --mount type=bind,source=$HOME/.ssh,target=/root/.ssh,readonly \
    --mount type=bind,source=$(go env GOCACHE),target=/gocache \
    --mount type=bind,source=$(go env GOMODCACHE),target=/gomodcache \
    thingify-go-build
```

## Alternative initial setup
As this repo depends on private GitHub repos, the following initial setup is required if not using the Docker build:

- Assuming SSH keys have been setup with access to the `github.com/thingify-app/` organisation, we want Go (or any other Git user) to default to SSH instead of HTTPS for `github.com` URLs. Modify `~/.gitconfig` to add the following lines:
  ```
  [url "ssh://git@github.com/"]
  	insteadOf = https://github.com/
  ```
- If the SSH key is password-protected, it should be added to the ssh-agent first:
  ```
  ssh-agent bash
  ssh-add ~/.ssh/<key_filename>
  ```
- Configure Go to treat all repos within `github.com/thingify-app/` as private. This just disables the use of the central Go package proxy:
  ```
  go env -w GOPRIVATE=github.com/thingify-app/*
  ```
