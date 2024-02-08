package widget

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type FocusableEntry struct {
	widget.Entry
	onFocusChanged func(bool)
}

func NewFocusableEntry() *FocusableEntry {
	f := &FocusableEntry{}
	f.ExtendBaseWidget(f)
	return f
}
func NewFocusableEntryWithData(data binding.String) *FocusableEntry {
	entry := NewFocusableEntry()
	entry.Bind(data)

	return entry
}

func (e *FocusableEntry) SetOnFocusChanged(listener func(bool)) {
	e.onFocusChanged = listener
}

func (e *FocusableEntry) FocusGained() {
	e.Entry.FocusGained()
	if e.onFocusChanged != nil {
		e.onFocusChanged(true)
	}
}

func (e *FocusableEntry) FocusLost() {
	e.Entry.FocusLost()
	if e.onFocusChanged != nil {
		e.onFocusChanged(false)
	}
}
