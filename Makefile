NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
UNAME := $(shell uname -s)
ifeq ($(UNAME),Darwin)
ECHO=echo
else
ECHO=/bin/echo -e
endif

all: deps
	$(ECHO) "$(OK_COLOR)==> Building plugin ...$(NO_COLOR)"
	@docker run -t -i --rm=true \
	--volumes-from gopath \
	-v `pwd`:/gopath/src/github.com/tonnydourado/packer-builder-docker-ssh \
	-w /gopath/src/github.com/tonnydourado/packer-builder-docker-ssh \
	-m=1g \
	google/golang \
	/bin/bash -c 'go build -o docker_ssh main.go'
	@sudo chown antonio:antonio docker_ssh

deps: gopath
	$(ECHO) "$(OK_COLOR)==> Downloading dependencies ...$(NO_COLOR)"
	@docker run -t -i --rm=true \
	--volumes-from gopath \
	-m=1g \
	google/golang \
	/bin/bash -c 'go get -u github.com/mitchellh/gox && go get github.com/mitchellh/packer && cd /gopath/src/github.com/mitchellh/packer && make deps'

test: deps
	@$(ECHO) "$(OK_COLOR)==> Testing plugin...$(NO_COLOR)"
	@docker run -t -i --rm=true \
	--volumes-from gopath \
	-v `pwd`:/gopath/src/github.com/tonnydourado/packer-builder-docker-ssh \
	-w /gopath/src/github.com/tonnydourado/packer-builder-docker-ssh \
	-m=1g \
	google/golang \
	/bin/bash -c 'go test ./builder/docker_ssh'

clean:
	@rm docker_ssh
	@docker ps --all | grep gopath | cut -b -12 | xargs docker rm

gopath:
	$(ECHO) "$(OK_COLOR)==> Creating data container 'gopath' ... $(NO_COLOR)"
	@test -z "`docker ps -a | grep gopath`" && \
	docker run -d --name gopath -v /gopath ubuntu:14.04 true || \
	$(ECHO) "$(OK_COLOR)==> 'gopath' volume container already exists$(NO_COLOR)"
