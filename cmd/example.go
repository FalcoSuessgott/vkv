package cmd

import (
	"fmt"

	"github.com/FalcoSuessgott/vkv/internal/convert"
	"github.com/FalcoSuessgott/vkv/pkg/example"
	"github.com/spf13/cobra"
)

const (
	numberOfArgs = 2
)

type exampleOptions struct {
	multiply bool
	add      bool
}

func defaultExampleOptions() *exampleOptions {
	return &exampleOptions{}
}

func newExampleCmd() *cobra.Command {
	o := defaultExampleOptions()

	cmd := &cobra.Command{
		Use:          "example",
		Short:        "example subcommand which adds or multiplies two given integers",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(numberOfArgs),
		RunE:         o.run,
	}

	cmd.Flags().BoolVarP(&o.multiply, "multiply", "m", o.multiply, "multiply")
	cmd.Flags().BoolVarP(&o.add, "add", "a", o.add, "add")

	return cmd
}

func (o *exampleOptions) run(cmd *cobra.Command, args []string) error {
	values, err := o.parseArgs(args)
	if err != nil {
		return err
	}

	if o.multiply {
		fmt.Fprintf(cmd.OutOrStdout(), "%d\n", example.Multiply(values[0], values[1]))
	}

	if o.add {
		fmt.Fprintf(cmd.OutOrStdout(), "%d\n", example.Add(values[0], values[1]))
	}

	return nil
}

func (o *exampleOptions) parseArgs(args []string) ([]int, error) {
	values := make([]int, 2) //nolint: gomnd

	for i, a := range args {
		v, err := convert.ToInteger(a)
		if err != nil {
			return nil, fmt.Errorf("error converting to integer: %w", err)
		}

		values[i] = v
	}

	return values, nil
}
