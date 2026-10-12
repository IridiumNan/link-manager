/*
Copyright © 2026 IridiumNan
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/IridiumNan/link-manager/internal/button"
	"github.com/IridiumNan/link-manager/internal/config"
	"github.com/IridiumNan/link-manager/internal/manager"
	"github.com/IridiumNan/link-manager/internal/models"
	"github.com/IridiumNan/link-manager/internal/teamodel"
	"github.com/IridiumNan/link-manager/pkgs/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "select then modify the link",
	Long:  `run command run tui for you to modify link efficiently`,
	Run:   execRun,
}

func execRun(cmd *cobra.Command, args []string) {
	// init logger
	logFile := models.SetLogger()

	slog.Debug("starting run command")

	defer logFile.Close()

	// load configuration file
	c, err := config.LoadConfig()

	if err == config.ErrInit {
		slog.Info(err.Error())
		return
	}
	if err != nil {

		slog.Error(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}

	choice, err := selectDir(utils.MapK(c.DirMap))
	if err != nil {
		slog.Error(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}

	if choice == "" {
		slog.Info("quit dir select without change")

		fmt.Print("bye")
		return
	}

	slog.Info("choose to modify link dir", "src_dir", choice, "dst_dir", c.DirMap[choice])

	b := button.NewBuilder(choice, c.DirMap[choice])
	b.SetNameMap(c.FileMap)

	err = modifyLinks(b)
	if err != nil {
		slog.Error(err.Error())
		fmt.Println(err)
		os.Exit(1)
	}

	slog.Debug("exiting run command without error")
}

func modifyLinks(b *button.ButtonBuilder) error {
	buttons, err := b.Build()
	if err != nil {
		return err
	}

	bMgr := manager.NewButtonMgr(buttons)

	lMgr := teamodel.NewLinkManager(bMgr)
	lMgr.SetHeader("modify below links")
	lMgr.SetFooter(`ctrl+c / q  to quit
		up / k to move cursor up
		down / j to move cursor down
		space  to modify status
		enter  to save changes`)

	p := tea.NewProgram(lMgr)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error while modifying link status, err: %s", err.Error())
	}

	fmt.Println(lMgr.Msg)

	if lMgr.Err != nil {
		fmt.Println(lMgr.Err)
	} else {
		color.Green("No error")
	}

	fmt.Println("\nbye")

	return nil
}

// selectDir request user select a specific dirs
func selectDir(dirs []string) (string, error) {
	dirSelector := teamodel.NewSelector(dirs)

	dirSelector.SetHeader("select the dir you want to probe or change link status")

	dirSelector.SetFooter(`ctrl+c / q  to quit
		up / k  to move cursor up
		down / j  to move cursor down
		enter  to choose and continue`)

	p := tea.NewProgram(dirSelector)

	if _, err := p.Run(); err != nil {
		return "", fmt.Errorf("error while running dir select: %s", err.Error())
	}

	return dirSelector.Choice(), nil
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
