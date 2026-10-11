package link

import (
	"fmt"
	"os"

	"github.com/IridiumNan/link-manager/internal/models"
	"github.com/fatih/color"
)

type LinkButton struct {
	linkSrc string

	linkDst string

	status models.Status
}

// TurnOn change button status to [models.StatusOn]
func (lb *LinkButton) TurnOn() {
	lb.status = models.StatusOn
}

// TurnOff change button status to [models.StatusOff]
func (lb *LinkButton) TurnOff() {
	lb.status = models.StatusOff
}

// Switch change reverse the button status
// if [models.StatusNone], it will not changed
func (lb *LinkButton) Switch() {
	switch lb.status {
	case models.StatusOff:
		lb.TurnOn()
	case models.StatusOn:
		lb.TurnOff()
	}
}

// String return a string name with color
func (lb *LinkButton) Name() string {
	return lb.strWithColor()
}

// create a string based on the status and link src, dst
// with format src -[]-> dst
// [] can be
// '-' for [models.StatusOn]
// 'x' for [models.StatusOff]
// '?' for [models.StatusNone]
//
// if colored, return the ---> with green, -x-> with red, -?-> with yellow
func (lb *LinkButton) str(colored bool) string {
	linkMarkBase := "-%s->"
	var linkMark string

	if colored {
		switch lb.status {
		case models.StatusOn:
			linkMark = color.GreenString(linkMarkBase, "-")
		case models.StatusOff:
			linkMark = color.RedString(linkMarkBase, "x")
		case models.StatusNone:
			linkMark = color.YellowString(linkMarkBase, "?")
		}
	} else {
		switch lb.status {
		case models.StatusOn:
			linkMark = fmt.Sprintf(linkMarkBase, "-")
		case models.StatusOff:
			linkMark = fmt.Sprintf(linkMarkBase, "x")
		case models.StatusNone:
			linkMark = fmt.Sprintf(linkMarkBase, "?")
		}
	}

	return fmt.Sprintf("%s %s %s", lb.linkSrc, linkMark, lb.linkDst)
}

// strWithouColor is alias to [LinkButton.str](false)
func (lb *LinkButton) strWithouColor() string {
	return lb.str(false)
}

// strWithColoris alias to [LinkButton.str](true)
func (lb *LinkButton) strWithColor() string {
	return lb.str(true)
}

// linkOnOperate create sym link if link not exist
// if link has existed, return [models.MsgNoChange]
func (lb *LinkButton) linkOnOperate() (models.ModifyMsg, error) {
	dst, err := os.Lstat(lb.linkDst)

	if os.IsNotExist(err) {

		err := os.Symlink(lb.linkSrc, lb.linkDst)
		if err != nil {
			return models.MsgFailed, fmt.Errorf("error while creating sym link %s -> %s, err: %s", lb.linkSrc, lb.linkDst, err.Error())
		}

		lb.TurnOn()
		return models.NewMsg("symlink: %s created", lb.strWithouColor()), nil
	}

	if !IsSymlink(dst.Mode()) {
		// it can be a dir, or file, DONT touch it
		return models.NewMsg("dst %s descriptor has exist, but not a link, it's mode: %s", lb.linkDst, dst.Mode()), os.ErrExist
	}
	return models.MsgNoChange, nil
}

// linkOffOperate try to remove the exist symlink
// if link not exist, it return [models.MsgNoChange] and nil
func (lb *LinkButton) linkOffOperate() (models.ModifyMsg, error) {
	dst, err := os.Lstat(lb.linkDst)

	if os.IsNotExist(err) {
		return models.MsgNoChange, nil
	}

	if !IsSymlink(dst.Mode()) {
		return models.NewMsg("dst %s descriptor has exist, but not a link, it's mode: %s, no change apply to it", lb.linkDst, dst.Mode()), os.ErrExist
	}

	err = os.Remove(lb.linkDst)
	if err != nil {
		return models.MsgFailed, fmt.Errorf("error while removing symlink: %s, err: %s", lb.linkDst, err.Error())
	}

	return models.NewMsg("symlink: %s has been removed", lb.strWithouColor()), nil
}

// Operate change the real system link based on current status
// if no change needed, it return [models.MsgNoChange]
// if status is [models.StatusNone], do nothing then return [models.MsgUnknown]
func (lb *LinkButton) Operate() (models.ModifyMsg, error) {
	switch lb.status {

	case models.StatusOn:
		return lb.linkOnOperate()

	case models.StatusOff:
		return lb.linkOffOperate()

	}
	return models.MsgUnknown, fmt.Errorf("unknown link status: %s", lb.Name())
}
