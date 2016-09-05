package configurer

import (
	"fmt"
	"os"
	"os/exec"

	"code.cloudfoundry.org/guardian/kawasaki/netns"
)

type ContainerConfigurer struct {
	pid    int
	execer *netns.Execer
}

func New(pid int) *ContainerConfigurer {
	return &ContainerConfigurer{
		pid:    pid,
		execer: &netns.Execer{},
	}
}

func (c *ContainerConfigurer) Exec(cmd *exec.Cmd) error {
	netnsFile, err := os.Open(fmt.Sprintf("/proc/%d/ns/net", c.pid))
	defer netnsFile.Close()
	if err != nil {
		return err
	}

	cbFunc := func() error {
		return cmd.Run()
	}

	return c.execer.Exec(netnsFile, cbFunc)
}
