/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/IridiumNan/link-manager/internal/config"
	"github.com/IridiumNan/link-manager/internal/models"
	"github.com/IridiumNan/link-manager/internal/teamodel"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: execRun,
}

func execRun(cmd *cobra.Command, args []string) {
	// init logger
	logFile := models.SetLogger()

	defer logFile.Close()

	// load configuration file
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error while loading configuration file: %s", err.Error())
	}

	dirs := make([]string, 0, len(c.DirMap))

	for dir := range c.DirMap {
		dirs = append(dirs, dir)
	}

	choice := selectDir(dirs)

	fmt.Printf("your choice is %s -> %s", choice, c.DirMap[choice])
}

func selectDir(dirs []string) string {
	dirSelector := teamodel.NewSelector(dirs)

	p := tea.NewProgram(dirSelector)

	if _, err := p.Run(); err != nil {
		log.Fatalf("error while running dir select: %s", err.Error())
	}

	return dirSelector.Choice()
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
