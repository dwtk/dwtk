package cli

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type Program struct {
	Version     string
	Description string
	SetFlags    func(fs *flag.FlagSet) error
	Pre         func(args []string) error
	Post        func(args []string) error

	name     string
	fs       *flag.FlagSet
	commands []*Command
}

type Command struct {
	Name        string
	Usage       string
	Description string
	SetFlags    func(fs *flag.FlagSet) error
	Run         func(args []string) error
}

var (
	version bool
	help    bool
)

func (p *Program) AddCommand(command *Command) {
	p.commands = append(p.commands, command)
}

func (p *Program) usageFlags() {
	fmt.Fprintf(p.fs.Output(), "Flags:\n")
	p.fs.PrintDefaults()
}

func (p *Program) usageCommand(command *Command) {
	if command.Description != "" {
		fmt.Fprintf(p.fs.Output(), "%s\n\n", command.Description)
	}

	fmt.Fprintf(p.fs.Output(), "Usage:\n  %s %s [flags]", p.name, command.Name)
	if command.Usage != "" {
		fmt.Fprintf(p.fs.Output(), " %s", command.Usage)
	}
	fmt.Fprintf(p.fs.Output(), "\n\n")

	p.usageFlags()
}

func (p *Program) usage() {
	if p.Description != "" {
		fmt.Fprintf(p.fs.Output(), "%s\n\n", p.Description)
	}

	fmt.Fprintf(p.fs.Output(), "Usage:\n  %s", p.name)
	if len(p.commands) > 0 {
		fmt.Fprintf(p.fs.Output(), " [command]")
	}
	fmt.Fprintf(p.fs.Output(), "\n\n")

	if len(p.commands) > 0 {
		fmt.Fprintf(p.fs.Output(), "Commands:\n")
		for _, cmd := range p.commands {
			fmt.Fprintf(p.fs.Output(), "  %s\n", cmd.Name)
			if cmd.Description != "" {
				fmt.Fprintf(p.fs.Output(), "    \t%s\n", cmd.Description)
			}
		}
		fmt.Fprintf(p.fs.Output(), "\n")
	}

	p.usageFlags()
}

func (p *Program) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("argument list must include binary name")
	}

	p.name = path.Base(args[0])
	if p.name == "" {
		return fmt.Errorf("binary name not found")
	}

	p.fs = flag.NewFlagSet(p.name, flag.ExitOnError)
	p.fs.BoolVar(&version, "v", false, "show version and exit")
	p.fs.BoolVar(&help, "h", false, "show help and exit")
	if p.SetFlags != nil {
		if err := p.SetFlags(p.fs); err != nil {
			return err
		}
	}
	p.fs.Usage = p.usage

	var (
		command *Command
		start   = 1
	)
	if len(args) > 1 {
		for _, cmd := range p.commands {
			if cmd.Name != args[1] {
				continue
			}

			if cmd.SetFlags != nil {
				if err := cmd.SetFlags(p.fs); err != nil {
					return err
				}
			}

			start = 2
			p.fs.Init(p.name+" "+cmd.Name, flag.ExitOnError)
			p.fs.Usage = func() {
				p.usageCommand(cmd)
			}
			command = cmd

			break
		}
	}

	if err := p.fs.Parse(args[start:]); err != nil {
		return err
	}

	if version {
		fmt.Println(p.name, p.Version)
		return nil
	}

	if command == nil || command.Run == nil || help {
		p.fs.Usage()
		os.Exit(2)
	}

	errs := []string{}
	if p.Pre != nil {
		if err := p.Pre(p.fs.Args()); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) == 0 {
		if err := command.Run(p.fs.Args()); err != nil {
			errs = append(errs, err.Error())
		}

		if p.Post != nil {
			if err := p.Post(p.fs.Args()); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	err := func() error {
		switch len(errs) {
		case 0:
			return nil
		case 1:
			return fmt.Errorf("Error: %s", errs[0])
		default:
			return fmt.Errorf("Errors:\n  - %s", strings.Join(errs, "\n  - "))
		}
	}()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
