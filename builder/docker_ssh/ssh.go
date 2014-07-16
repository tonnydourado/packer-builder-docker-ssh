package docker_ssh

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"log"
	"os/exec"
)

// SSHAddress returns a function that can be given to the SSH communicator
// for determining the SSH address
func SSHAddress(port int) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		containerId := state.Get("container_id").(string)
		host, err := exec.Command("docker", "inspect", "--format", "{{ .NetworkSettings.IPAddress }}", containerId).Output()
		if err != nil {
			log.Fatal(err)
		}
		return fmt.Sprintf("%s:%d", host, port), nil
	}
}
