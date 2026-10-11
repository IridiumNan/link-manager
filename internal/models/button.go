package models

import (
	"fmt"
	"io"
)

type Status int

const (
	// StatusNone means this button's status is unknown
	// Just use it when initialize a new button
	StatusNone = iota

	StatusOn
	StatusOff
)

type Button interface {
	// Switch change button's status
	Switch()

	// TurnOn change status to [StatusOn]
	TurnOn()

	// TurnOff change status to [StatusOff]
	TurnOff()

	// Status return current button status
	// It can be [StatusNone] [StatusOn] [StatusOff]
	Status() Status
}

type Name interface {
	// String return name without '\n'
	Name() string
}

// ModifyMsg should not end with '\n'
type ModifyMsg string

const (
	MsgNoChange ModifyMsg = "not change"
	MsgFailed   ModifyMsg = "fail to make change"
	MsgUnknown  ModifyMsg = "unknown link status"
)

func NewMsg(format string, a ...any) ModifyMsg {
	return ModifyMsg(fmt.Sprintf(format, a...))
}

type Operator interface {
	// Make real change based on it's current status
	//
	// if status not changed, it will return [MsgNoChange], nil
	Operate() (ModifyMsg, error)
}

type NameButtonOperator interface {
	// visit name by [Name.String]
	Name

	// change status by [Button.Switch] [Button.TurnOn] [Button.TurnOff]
	// check status by [Button.Status]
	Button

	// do change by [Operator.Operate]
	Operator
}

type NameCloser interface {
	Name

	io.Closer
}
