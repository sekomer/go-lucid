package main

import (
	"fmt"
	bootnode "go-lucid/cmd/bootnode"
	devnode "go-lucid/cmd/devnode"
	fullnode "go-lucid/cmd/fullnode"
	"go-lucid/node"
	"io"
	"log"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/rand"
	"gopkg.in/yaml.v3"
)

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))

	log.Default().SetFlags(log.LstdFlags | log.Lshortfile)

	var yamlFile string
	var boot bool
	var dev bool

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version-flag",
		Aliases: []string{"V", "v", "version"},
		Usage:   "print the version",
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help-flag",
		Aliases: []string{"help", "h"},
		Usage:   "show help",
	}

	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "gucid",
		Usage:                "go execution layer implementation of the lucid protocol",
		Version:              "0.1.0",
		AllowExtFlags:        false,
		Suggest:              true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				TakesFile:   true,
				Value:       "config/fullnode.yaml",
				Usage:       "Load configuration from `FILE`",
				Destination: &yamlFile,
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 && ctx.NumFlags() == 0 {
				cli.ShowAppHelp(ctx)
				return nil
			}

			return nil

		},
		Commands: []*cli.Command{
			{
				Name:    "node",
				Usage:   "manage lucid nodes",
				Aliases: []string{"n"},
				Subcommands: []*cli.Command{
					{
						Name:    "start",
						Usage:   "start the lucid node",
						Aliases: []string{"s"},
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:        "boot",
								Usage:       "start the bootnode",
								Value:       false,
								Destination: &boot,
							},
							&cli.BoolFlag{
								Name:        "dev",
								Usage:       "start the dev",
								Value:       false,
								Destination: &dev,
							},
						},
						Action: func(*cli.Context) error {
							log.Println("yamlFile:", yamlFile)

							yamlReader, err := os.Open(yamlFile)
							if err != nil {
								return err
							}
							defer yamlReader.Close()

							buf, err := io.ReadAll(yamlReader)
							if err != nil {
								return err
							}

							c := &node.FullNodeConfig{}
							err = yaml.Unmarshal(buf, c)
							if err != nil {
								return fmt.Errorf("in file %q: %w", yamlFile, err)
							}

							if boot {
								bootnode.Start(c)
							} else if dev {
								devnode.Start(c)
							} else {
								fullnode.Start(c)
							}

							return nil
						},
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
