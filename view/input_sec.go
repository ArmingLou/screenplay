package view

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"screenplay/controller"
	"strconv"
)

func CreatInputSec(controller *controller.InputSecController) *widget.Entry {

	inputSec := widget.NewEntryWithData(controller.GetBinding())
	inputSec.OnChanged = func(s string) {
		if controller.OnChange != nil {
			i, e := strconv.Atoi(s)
			if e != nil {
				return
			}
			if i < 1 {
				return
			}
			controller.OnChange(i)
		}
	}
	inputSec.SetPlaceHolder("Enter seconds")
	inputSec.Validator = func(s string) error {
		i, e := strconv.Atoi(s)
		if e != nil {
			return e
		}
		if i < 1 {
			return fmt.Errorf("must > 0")
		}
		return nil
	}
	return inputSec
}
