package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "flutter-go-bridge",
		Short: "FFI bridge generator between Dart/Flutter and Go",
	}

	rootCmd.AddCommand(cmdGenerate)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
