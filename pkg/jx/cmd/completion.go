package cmd

import (
	"io"

	"github.com/jenkins-x/jx/pkg/jx/cmd/templates"
	cmdutil "github.com/jenkins-x/jx/pkg/jx/cmd/util"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

const boilerPlate = ""

var (
	completion_long = templates.LongDesc(`
		Output shell completion code for the given shell (bash or zsh).

		This command prints shell code which must be evaluation to provide interactive
		completion of jx commands.

		    $ source <(jx completion bash)

		will load the jx completion code for bash. Note that this depends on the
		bash-completion framework. It must be sourced before sourcing the jx
		completion, e.g. on the Mac:

		    $ brew install bash-completion
		    $ source $(brew --prefix)/etc/bash_completion
		    $ source <(jx completion bash)

		On a Mac it often works better to generate a file with the completion and source that:

			$ jx completion bash > ~/.jx/bash
			$ source ~/.jx/bash

		If you use zsh[1], doing the following will enable jx zsh completion.

		With oh-my-zsh:

			Create a custom plugin for jx since none is currently available:

			    $ mkdir -p ${ZSH_CUSTOM}/plugins/jx
			    $ jx completion zsh>${ZSH_CUSTOM}/plugins/jx/_jx

			In your .zshrc file, enable the plugin by adding it to your plugin list:

			    plugins=(… jx)

			Finally, start a new shell to have jx zsh completions working.

		Manual installation:

			First create a directory where you will put the completion script:

			    $ mkdir ${HOME}/.jx-completion

			Then add the following to your .zshrc file:

			    fpath=(${HOME}/.jx-completion $fpath)
			    jx completion zsh>${HOME}/.jx-completion/_jx

			Finally, start a new shell to have jx zsh completions working.

		[1] zsh completions are only supported in versions of zsh >= 5.2`)
)

var (
	completion_shells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

func NewCmdCompletion(f cmdutil.Factory, out io.Writer) *cobra.Command {
	options := &CommonOptions{
		Factory: f,
		Out:     out,
		Err:     out,
	}

	shells := []string{}
	for s := range completion_shells {
		shells = append(shells, s)
	}

	cmd := &cobra.Command{
		Use:   "completion SHELL",
		Short: "Output shell completion code for the given shell (bash or zsh)",
		Long:  completion_long,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			cmdutil.CheckErr(err)
		},
		ValidArgs: shells,
	}

	return cmd
}

func (o *CommonOptions) Run() error {
	shells := []string{}
	for s := range completion_shells {
		shells = append(shells, s)
	}
	var ShellName string
	cmd := o.Cmd
	args := o.Args
	if len(args) == 0 {
		prompts := &survey.Select{
			Message:  "Shell",
			Options:  shells,
			PageSize: len(shells),
			Help:     "The name of the shell",
		}
		err := survey.AskOne(prompts, &ShellName, nil)
		if err != nil {
			return err
		}

	}
	if len(args) > 1 {
		return cmdutil.UsageError(cmd, "Too many arguments. Expected only the shell type.")
	}
	if ShellName == "" {
		ShellName = args[0]
	}

	run, found := completion_shells[ShellName]

	if !found {
		return cmdutil.UsageError(cmd, "Unsupported shell type %q.", args[0])
	}

	return run(o.Out, cmd.Parent())
}

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	if boilerPlate != "" {
		_, err := out.Write([]byte(boilerPlate))
		if err != nil {
			return err
		}
	}
	return cmd.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, cmd *cobra.Command) error {
	if boilerPlate != "" {
		_, err := out.Write([]byte(boilerPlate))
		if err != nil {
			return err
		}
	}
	return cmd.GenZshCompletion(out)
}
