package main

import (
	c "github.com/gookit/color"
	"github.com/jirallreadyforthis/cmd"
	"os"
)

func main() {
	cobraCmd, err := cmd.Make()
	if err != nil {
		c.Printf("<red>jirallreadyforthis: building cmd</> %v", err)
		os.Exit(1)
	}

	if err := cobraCmd.Execute(); err != nil {
		c.Printf("<red>jirallreadyforthis:</> %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
