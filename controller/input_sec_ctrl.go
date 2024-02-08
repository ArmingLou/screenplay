package controller

import (
	"fyne.io/fyne/v2/data/binding"
	"strconv"
)

type InputSecController struct {
	inputSecStr binding.String
	OnChange    func(sec int)
}

func NewInputSecController() *InputSecController {
	ib := binding.NewString()
	ib.Set("1")
	return &InputSecController{
		inputSecStr: ib,
	}
}

func (c *InputSecController) SetValue(sec int) {
	c.inputSecStr.Set(strconv.Itoa(sec))
}

func (c *InputSecController) GetValue() (sec int, err error) {
	s, e := c.inputSecStr.Get()
	if e != nil {
		return 0, e
	}
	i, e := strconv.Atoi(s)
	if e != nil {
		return 0, e
	}
	return i, nil
}

func (c *InputSecController) GetBinding() binding.String {
	return c.inputSecStr
}
