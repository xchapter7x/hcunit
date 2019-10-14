package commands

import "fmt"

type VersionCommand struct {
	Version   string
	Buildtime string
	Platform  string
}

func (s *VersionCommand) Execute(args []string) error {
	fmt.Printf("ver: %v date: %v platform: %v \n", s.Version, s.Buildtime, s.Platform)
	return nil
}
