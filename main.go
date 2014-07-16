package main

import (
    "github.com/tonnydourado/packer-builder-docker-ssh/builder/docker_ssh"
    "github.com/mitchellh/packer/packer/plugin"
)

func main() {
    server, err := plugin.Server()
    if err != nil {
        panic(err)
    }
    server.RegisterBuilder(new(docker_ssh.Builder))
    server.Serve()
}
