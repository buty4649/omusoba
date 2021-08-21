package cmd

import (
	"os"

	"buty4649/omusoba/version"

	"github.com/spf13/cobra"
)

var (
	cfgUsb    string
	cfgFormat string

	rootCmd = &cobra.Command{
		Use:          "omusoba",
		Version:      version.Version,
		Short:        "omusoba is a CLI tool for OMRON 2JCIE-BL/BU",
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgUsb, "usb", "u", "", "USB device path")
	//rootCmd.MarkPersistentFlagRequired("usb")

	rootCmd.PersistentFlags().StringVarP(&cfgFormat, "format", "f", "label", "the output format label/yaml/json")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	return nil
}
