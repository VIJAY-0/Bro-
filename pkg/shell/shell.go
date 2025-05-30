package shell

import (
	"fmt"
	"os/exec"
)

func Run_command(command string) (string, error) {

	fmt.Printf("executing comamand on shell < %s >\n", command)
	cmd := exec.Command("bash", "-c", command)
	response, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}
	return string(response), nil

}
