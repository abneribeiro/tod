package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "tod",
    Short: "A simple CLI for compressing and decompressing files",
}

// Execute executes the root command.
func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize()
    rootCmd.PersistentFlags().StringP("path", "p", "", "path to the file or directory")
}
