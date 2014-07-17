# packer-builder-docker-ssh


Packer builder that provisions containers using SSH, instead of the docker attach command, to use with images like stackbrew/ubuntu-upstart, that have a proper init process.

## Install

You will need at least [Docker](docker.io) and make to build this plugin. Once they are installed, clone this repository:

```
$ git clone git@github.com:tonnydourado/packer-builder-docker-ssh.git
```

Then just run `make`. This will create a container called gopath, with a [data volume](https://docs.docker.com/userguide/dockervolumes/) containing the `GOPATH` in which the plugin will be built. All other targets (like `test`) will run inside containers with limited memory, sharing this volume. It might sound complicated, but I found this method cleanner than installing Go in my host system.

You might speed up the first build by pulling the `google/golang` and `ubuntu:14.04` Docker images before it.

## Basic Example

Here's a basic example of how to use this builder:

```
{
  "builders": [
    {
      "type": "docker_ssh",
      "image": "stackbrew/ubuntu-upstart:14.04",
      "export_path": "gedvic_base.tar",
      "ssh_username": "root",
      "ssh_password": "docker.io"

    }
  ],

  "provisioners": [
    {
      "type": "shell",
      "inline": "apt-get update --fix-missing"
    }
  ]
}
```

Instead of password, an identity key can be used. It's also possible to specify a alternative port to be used (the default  22)

## Configuration Reference

The reference of available configuration options is listed below.

### Required parameters:

 * `export_path` (string) - The path where the final container will be exported as a tar file..
 * `image` (string) - The base image for the Docker container that will be started. This image will be pulled from the Docker registry if it doesn't already exist.
 * `ssh_username` (string) - The name of the resulting image that will appear in your account. This must be unique. To help make this unique, use a function like timestamp.
 * `ssh_password` (string) - The password to be used for the ssh connection. Cannot be combined with `ssh_private_key_file`.
 * `ssh_private_key_file` (string) - The filename of the ssh private key to be used for the ssh connection. E.g. `/home/user/.ssh/identity_rsa`.

### Optional parameters:

 * `pull` (boolean) - If true, the configured image will be pulled using docker pull prior to use. Otherwise, it is assumed the image already exists and can be used. This defaults to true if not set.
 * `run_command` (array of strings) - An array of arguments to pass to `docker` in order to run the container. By default this is set to `["run", "-d", "-v", "{{.Volumes}}", "{{.Image}}", "/sbin/init"]`. As you can see, you have a couple template variables to customize, as well.

## Contribute

Contributions are always welcome! Pull requests, issues, all that stuff, go alredy and do it =)

### TODO
* Configure travis CI or any alternative to automatically test and build the code
* Build and package it
* Create a project page, or something like that
