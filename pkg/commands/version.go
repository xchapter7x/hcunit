package commands

import "fmt"

type VersionCommand struct {
	Version   string
	Buildtime string
	Platform  string
}

func (x *VersionCommand) Execute(args []string) error {
	fmt.Printf("ver: %v date: %v platform: %v \n", x.Version, x.Buildtime, x.Platform)
	return nil
}
