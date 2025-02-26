package main

import (
	"github.com/spf13/cobra"

	"github.com/shipa-corp/ketch/internal/deploy"
)

const (
	appDeployHelp = `
Roll out a new version of an application with an image.

Deploy from source code. <source> is path to source code. The image in this case is required
and will be built using the selected source code and platform and will be used to deploy the app. 
  ketch app deploy <app name> <source> -i myregistry/myimage:latest 

  Ketch looks for Procfile and ketch.yaml inside the source directory by default 
  but you can provide a custom path with --procfile or --ketch-yaml.

Deploy from an image:
  ketch app deploy <app name> -i myregistry/myimage:latest

  Ketch uses the image's cmd and entrypoint but you can redefine what exactly to run with --procfile. 

`
)

// NewCommand creates a command that will run the app deploy
func newAppDeployCmd(params *deploy.Services) *cobra.Command {
	var options deploy.Options

	cmd := &cobra.Command{
		Use:   "deploy APPNAME [SOURCE DIRECTORY]",
		Short: "Deploy an app.",
		Long:  appDeployHelp,
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.AppName = args[0]
			if len(args) == 2 {
				options.AppSourcePath = args[1]
			}
			return deploy.New(options.GetChangeSet(cmd.Flags())).Run(cmd.Context(), params)
		},
	}

	cmd.Flags().StringVarP(&options.Image, deploy.FlagImage, deploy.FlagImageShort, "", "Name of the image to be deployed.")
	cmd.Flags().StringVar(&options.KetchYamlFileName, deploy.FlagKetchYaml, "", "Path to ketch.yaml.")

	cmd.Flags().StringVar(&options.ProcfileFileName, deploy.FlagProcFile, "", "Path to procfile.")
	cmd.Flags().BoolVar(&options.StrictKetchYamlDecoding, deploy.FlagStrict, false, "Enforces strict decoding of ketch.yaml.")
	cmd.Flags().IntVar(&options.Steps, deploy.FlagSteps, 2, "Number of steps for a canary deployment.")
	cmd.Flags().StringVar(&options.StepTimeInterval, deploy.FlagStepInterval, "", "Time interval between canary deployment steps. Supported min: m, hour:h, second:s. ex. 1m, 60s, 1h.")
	cmd.Flags().BoolVar(&options.Wait, deploy.FlagWait, false, "If true blocks until deploy completes or a timeout occurs.")
	cmd.Flags().StringVar(&options.Timeout, deploy.FlagTimeout, "20s", "Defines the length of time to block waiting for deployment completion. Supported min: m, hour:h, second:s. ex. 1m, 60s, 1h.")
	cmd.Flags().StringSliceVar(&options.SubPaths, deploy.FlagIncludeDirs, []string{"."}, "Optionally include additional source paths. Additional paths must be relative to source-path.")

	cmd.Flags().StringVarP(&options.Platform, deploy.FlagPlatform, deploy.FlagPlatformShort, "", "Platform name.")
	cmd.Flags().StringVarP(&options.Description, deploy.FlagDescription, deploy.FlagDescriptionShort, "", "App description.")
	cmd.Flags().StringSliceVarP(&options.Envs, deploy.FlagEnvironment, deploy.FlagEnvironmentShort, []string{}, "App env variables.")
	cmd.Flags().StringVarP(&options.Pool, deploy.FlagPool, deploy.FlagPoolShort, "", "Pool to deploy your app.")
	cmd.Flags().StringVarP(&options.DockerRegistrySecret, deploy.FlagRegistrySecret, "", "", "A name of a Secret with docker credentials. This secret must be created in the same namespace of the pool.")

	return cmd
}
