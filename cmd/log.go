/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"

	"github.com/IridiumNan/link-manager/internal/models"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Open log file with less command",
	Long: `This command will open your log file with less, 
	you should ues it when something wrong happened`,
	Run: execLog,
}

func execLog(cmd *cobra.Command, args []string) {
	logFile := models.SetLogger()

	slog.Debug("starting log command")

	defer logFile.Close()

	fmt.Println("the logFile path -> ", logFile.Name())

	// default viewer is less command
	lessCmd := exec.Command("less", logFile.Name())

	lessCmd.Stdin = os.Stdin
	lessCmd.Stdout = os.Stdout
	lessCmd.Stderr = os.Stderr

	if err := lessCmd.Run(); err != nil {
		logFile.Close()
		log.Fatal(err)
	}

	slog.Debug("log command exit without error")
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
