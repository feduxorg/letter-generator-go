package main

import (
	"os"
	"path/filepath"

	"github.com/feduxorg/letter-generator/letter_generator"
	lgos "github.com/feduxorg/letter-generator/os"
	"github.com/feduxorg/letter-generator/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// The current build version.
var AppVersionNumber = "dev"

// SHA-value of git commit
var CommitHash = "HEAD"

// Date of build
var BuildDate = "NOW"

type Cli struct{}

func (p *Cli) Run(args []string) error {
	app := cli.NewApp()
	app.Name = "letter-generator"
	app.Version = AppVersionNumber + "-" + CommitHash + "-" + BuildDate

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose,V",
			Usage: "activate verbose logging",
		},
		cli.BoolFlag{
			Name:  "show-config,C",
			Usage: "Show configuration",
		},
	}

	app.Action = func(c *cli.Context) error {
		var workDir string

		if c.Args().Get(0) != "" {
			workDir = c.Args().Get(0)
		} else {
			workDir = getCwd()
		}

		config := buildConfig(workDir)
		parseGlobalOptions(c, config)

		err := build(config)

		if err != nil {
			return err
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "initialize current directory",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose, V",
					Usage: "activate verbose logging",
				},
				cli.StringFlag{
					Name:  "template-source, t",
					Usage: "git repository with templates",
				},
			},
			Action: func(c *cli.Context) error {
				var workDir string

				if c.Args().Get(0) != "" {
					workDir = c.Args().Get(0)
				} else {
					workDir = getCwd()
				}

				config := buildConfig(workDir)
				parseGlobalOptions(c, config)

				if c.String("template-source") != "" {
					config.TemplateSource = c.String("template-source")
				}

				err := initialize(workDir, config)

				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "build letters based on information in current directory",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "verbose,V",
					Usage: "activate verbose logging",
				},
			},
			Action: func(c *cli.Context) error {
				var workDir string

				if c.Args().Get(0) != "" {
					workDir = c.Args().Get(0)
				} else {
					workDir = getCwd()
				}

				config := buildConfig(workDir)
				parseGlobalOptions(c, config)

				err := build(config)

				if err != nil {
					return err
				}

				return nil
			},
		},
	}

	app.Run(os.Args)

	return nil
}

func build(config letter_generator.Config) error {
	builder := api.LetterBuilder{}
	err := builder.Build(config)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func initialize(dir string, config letter_generator.Config) error {
	initializer := api.Initializer{}
	err := initializer.Init(dir, config)

	if err != nil {
		return err
	}

	return nil
}

func parseGlobalOptions(c *cli.Context, config letter_generator.Config) {
	if c.Bool("verbose") == true {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.WithFields(log.Fields{
		"verbose": c.Bool("verbose"),
	}).Debug("Parsing commandline options")
}

func getCwd() string {
	currentDir, err := os.Getwd()

	if err != nil {
		log.WithFields(log.Fields{
			"msg":    err.Error(),
			"status": "failure",
		}).Fatal("Getting current directory")
	}

	log.WithFields(log.Fields{
		"path":   currentDir,
		"status": "success",
	}).Debug("Getting current directory")

	return currentDir
}

func buildConfig(workDir string) letter_generator.Config {
	homeDir, err := lgos.HomeDirectory()

	if err != nil {
		log.WithFields(log.Fields{
			"msg":    err.Error(),
			"status": "failure",
		}).Fatal("Getting home directory of current user")
	}

	log.WithFields(log.Fields{
		"path":   homeDir,
		"status": "success",
	}).Debug("Getting home directory of current user")

	config := letter_generator.Config{}
	config.TemplateSource = filepath.Join(homeDir, ".local/share/letter-template/.git")
	config.ConfigDirectory = ".lg"
	config.RecipientsFile = filepath.Join(workDir, config.ConfigDirectory, "data/to.yaml")
	config.MetadataFile = filepath.Join(workDir, config.ConfigDirectory, "data/metadata.yaml")
	config.SenderFile = filepath.Join(workDir, config.ConfigDirectory, "data/from.yaml")
	config.TemplateFile = filepath.Join(workDir, config.ConfigDirectory, "templates/letter.tex.tt")
	config.AssetsDirectory = filepath.Join(workDir, config.ConfigDirectory, "assets")

	return config
}
