package prerun

import "github.com/spf13/cobra"

type CobraPreRun = func(cmd *cobra.Command, args []string) error

func Combine(preruns ...CobraPreRun) CobraPreRun {
	return func(cmd *cobra.Command, args []string) error {
		for _, f := range preruns {
			err := f(cmd, args)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
