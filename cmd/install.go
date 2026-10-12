/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/IridiumNan/link-manager/internal/button"
	"github.com/IridiumNan/link-manager/internal/config"
	"github.com/IridiumNan/link-manager/internal/models"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install all link loaded from config",
	Long:  `Install command will load configuration file then create all soft links`,
	Run:   execInstall,
}

func execInstall(cmd *cobra.Command, args []string) {
	logFile := models.SetLogger()
	slog.Debug("starting instal")

	defer logFile.Close()

	c, err := config.LoadConfig()
	if err != nil {
		slog.Error(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}

	for src, dst := range c.DirMap {
		fmt.Printf("\ninstalling for %s -> %s", src, dst)
		err := installOnce(src, dst, c.FileMap)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func installOnce(src string, dst string, nameMap map[string]string) error {
	var wg sync.WaitGroup

	b := button.NewBuilder(src, dst)

	b.SetNameMap(nameMap)

	buttons, err := b.Build()
	if err != nil {
		return fmt.Errorf("error while build buttons: %s", err.Error())
	}

	for _, btn := range buttons {
		wg.Go(func() {
			btn.TurnOn()
			msg, err := btn.Operate()
			if err != nil {
				slog.Error(err.Error())
			}
			fmt.Println(msg)
		})
	}

	wg.Wait()

	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
