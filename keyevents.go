package main

import "github.com/lxn/walk"

func ctrlTab(key walk.Key) {
	if key == walk.KeyTab && walk.ControlDown() {
		max := tabWidget.Pages().Len() - 1
		if max < 1 {
			return
		}
		index := tabWidget.CurrentIndex()
		if walk.ShiftDown() {
			index -= 1
			if index < 0 {
				index = max
			}
		} else {
			index += 1
			if index > max {
				index = 0
			}
		}
		tabWidget.SetCurrentIndex(index)
	}

	if key == walk.KeyEscape && walk.ShiftDown() {
		mw.ToggleBorder()
	}

	if key == walk.KeyF2 {
		mw.SetTransparency(-16)
	}
	if key == walk.KeyF4 {
		mw.SetTransparency(16)
	}
	if key == walk.KeyF3 {
		mw.ToggleTransparency()
	}
}

func insertCharacter(key walk.Key) rune {
	if walk.ControlDown() {
		switch key {
		case walk.KeyK:
			return fmtColor
		case walk.KeyB:
			return fmtBold
		case walk.KeyI:
			return fmtItalic
		case walk.KeyU:
			return fmtUnderline
		case walk.KeyS:
			return fmtStrikethrough
		case walk.Key0, walk.KeyNumpad0:
			return fmtReset
		}
	}
	return 0
}
