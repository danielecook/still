package commands

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/danielecook/still/src/schema"
	"github.com/danielecook/still/src/validate"
)

// Run - Entrypoint
func Run(Version string) {

	app := cli.NewApp()
	app.Name = "Still"
	app.Usage = "Validate CSV, TSV, and Excel data"
	app.Version = Version

	app.Authors = []*cli.Author{
		{
			Name:  "Daniel Cook",
			Email: "danielecook@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "schema",
			Aliases: []string{"s"},
			Usage:   "Load schema from `FILE`",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "validate",
			Usage: "validate <input> <schema>",
			Action: func(c *cli.Context) error {
				schemaFname := c.Args().Get(0)
				input := c.Args().Get(1)
				useSchema := schema.ParseSchema(schemaFname)
				// Rev direction of params
				result := validate.RunValidation(useSchema, input)
				if result == false {
					os.Exit(1)
				} else {
					os.Exit(0)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
