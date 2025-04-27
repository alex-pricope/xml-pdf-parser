package main

import (
	"github.com/alex-pricope/form-parser/cmd"
	"github.com/alex-pricope/form-parser/logging"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	// Decided not to inject the logger and use it globally like this to simplify the app
	logging.BoostrapLogger()

	var rootCmd = &cobra.Command{
		Use:     "parser",
		Short:   "Simple file parser",
		Example: "parser --file=input_file --sub=submission_file --from=xml --to=pdf --out=./output/",
		Run:     cmd.ParseCommand,
	}

	// Add the flags - can be extended with others
	rootCmd.Flags().StringP("file", "f", "", "file to parse")
	err := rootCmd.MarkFlagRequired("file")
	if err != nil {
		logging.Log.Error(err)
		return
	}

	rootCmd.Flags().StringP("sub", "s", "", "file to parse")
	err = rootCmd.MarkFlagRequired("sub")
	if err != nil {
		logging.Log.Error(err)
		return
	}

	rootCmd.Flags().String("from", "", "Input file type")
	err = rootCmd.MarkFlagRequired("from")
	if err != nil {
		logging.Log.Error(err)
		return
	}

	rootCmd.Flags().String("to", "", "Target file type")
	err = rootCmd.MarkFlagRequired("to")
	if err != nil {
		logging.Log.Error(err)
		return
	}

	rootCmd.Flags().StringP("out", "o", "", "Output folder")

	if err = rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
