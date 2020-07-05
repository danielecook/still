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
				useSchema.OutputOrder = c.String("order")
				// Rev direction of params
				result := validate.RunValidation(useSchema, input)
				if result == false {
					os.Exit(1)
				} else {
					os.Exit(0)
				}
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "yaml",
					Usage:    "Specify YAML Provider",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "order",
					Aliases:  []string{"o"},
					Value:    "data",
					Usage:    "Order output by data or schema",
					Required: false,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
