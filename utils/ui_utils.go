package utils

import (
	"fmt"
	"github.com/briandowns/spinner"
	"sync/atomic"

	"github.com/jroimartin/gocui"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

var LogWin *gocui.View
var GlobalGui *gocui.Gui

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/

var spinnerChars = []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}

func UseGuiQ() bool {
	return GlobalGui != nil
}

func SetGui(g *gocui.Gui) {
	GlobalGui = g
}

var nextSpinnerPos int32
var activeSpinners int32

type UniSpinner struct {
	plainSp *spinner.Spinner
	guiSp   *GuiSpinner
}

func NewSpinner(prompt string) *UniSpinner {

	atomic.AddInt32(&activeSpinners, 1)
	sp := &UniSpinner{}
	useGui := GlobalGui != nil
	if useGui {
		sp.guiSp = NewGuiSpinner(prompt)
	} else {
		sp.plainSp = spinner.New(spinner.CharSets[0], 150*time.Millisecond)
		sp.plainSp.Prefix = prompt
		sp.plainSp.Start()
	}

	return sp
}

func (sp *UniSpinner) Stop() {
	nv := atomic.AddInt32(&activeSpinners, -1)
	if nv <= 0 {
		atomic.StoreInt32(&activeSpinners, 0)
		atomic.StoreInt32(&nextSpinnerPos, 0)
	}
	if sp.plainSp != nil {
		sp.plainSp.Stop()
	} else if sp.guiSp != nil {
		sp.guiSp.Stop()
	}
}

type GuiSpinner struct {
	spinnerActive AtomBool
	v             *gocui.View
	viewName      string
}

func NewGuiSpinner(prompt string) *GuiSpinner {

	sp := &GuiSpinner{}
	sp.spinnerActive.Set()

	as := atomic.AddInt32(&nextSpinnerPos, 1) - 1
	sp.viewName = fmt.Sprintf("waitView%d", as)

	GlobalGui.Update(func(gui *gocui.Gui) error {
		sp.ShowWaitingView(prompt, sp.viewName, int(as))
		return nil
	})

	go func() {
		ndx := 0
		for sp.spinnerActive.Get() {
			ndx++
			if ndx >= len(spinnerChars) {
				ndx = 0
			}

			if sp.v != nil {
				GlobalGui.Update(func(gui *gocui.Gui) error {

					sp.v.Clear()
					sp.v.SetCursor(2, 0)
					fmt.Fprintf(sp.v, " %s", prompt)
					sp.v.SetCursor(len(prompt)+1, 0)
					io.WriteString(sp.v, spinnerChars[ndx])
					return nil
				})
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	return sp
}

func (sp *GuiSpinner) Stop() {
	if !sp.spinnerActive.Get() {
		return
	}
	sp.spinnerActive.Reset()
	time.Sleep(time.Millisecond * 10)
	GlobalGui.Update(func(gui *gocui.Gui) error {
		sp.DelWaitingView(sp.viewName)
		return nil
	})

}

func (sp *GuiSpinner) ShowWaitingView(prompt string, viewName string, offset int) (*gocui.View, error) {
	maxX, _ := GlobalGui.Size()
	l := len(prompt) + 5
	var err error
	sp.v, err = GlobalGui.SetView(viewName, maxX-l, 3*offset, maxX-1, 3*offset+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
	}
	sp.v.Wrap = false

	sp.v.Editable = false
	sp.v.SetCursor(2, 0)
	fmt.Fprintf(sp.v, " %s", prompt)
	sp.v.SetCursor(len(prompt)+1, 0)

	time.Sleep(time.Millisecond * 50)
	GlobalGui.SetViewOnTop(viewName)
	return sp.v, nil
}

func (sp *GuiSpinner) DelWaitingView(viewName string) error {
	if err := GlobalGui.DeleteView(viewName); err != nil {
		return err
	}
	switchView(GlobalGui, LogWin)
	return nil
}

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/

func ShowConfView(g *gocui.Gui, prompt string, action func(*string)) (*gocui.View, error) {
	maxX, maxY := g.Size()

	l := len(prompt) + 5
	v, err := g.SetView("confView", maxX/2-l/2-1, maxY/2-1, maxX/2+l/2, maxY/2+1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

	}
	v.Wrap = false
	g.Cursor = true
	v.Editable = true
	v.SetCursor(2, 0)
	fmt.Fprintf(v, " %s", prompt)
	v.SetCursor(len(prompt)+2, 0)

	if err := g.SetKeybinding("confView", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			str := v.Buffer()
			str = strings.Trim(str, "\n")
			str = str[len(prompt)+1:]
			str = strings.Trim(str, " ")
			//logger.Debugf("*", "response was '%s'", str)
			time.Sleep(time.Millisecond * 50)
			action(&str)

			g.DeleteKeybinding("confView", gocui.KeyEnter, gocui.ModNone)
			DelConfView(g)

			return nil
		}); err != nil {
		log.Panicln(err)
	}

	time.Sleep(time.Millisecond * 50)
	if _, err := g.SetCurrentView("confView"); err != nil {
		return nil, err
	}

	g.SetViewOnTop("confView")
	return v, nil
}

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/

func ShowTempView(g *gocui.Gui, name string, x0, y0, x1, y1 int, wrap, cursor, edit bool) (*gocui.View, error) {

	v, err := g.SetView(name, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
	}

	v.Wrap = wrap
	g.Cursor = cursor
	v.Editable = edit

	time.Sleep(time.Millisecond * 50)
	g.SetViewOnTop(name)

	return v, nil
}

func DelTempView(g *gocui.Gui, name string) error {
	if err := g.DeleteView(name); err != nil {
		return err
	}
	switchView(g, LogWin)
	return nil
}

func DelConfView(g *gocui.Gui) error {
	return DelTempView(g, "confView")
}

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/

func HandleShutdowns(ctx *ContextWithCancel) *ContextWithCancel {
	shutdownSignalChan := make(chan os.Signal, 1)
	signal.Notify(shutdownSignalChan, os.Interrupt, os.Kill)
	go func(ctx *ContextWithCancel) {
		for {
			ctx.Logger.Infof("*", "'%v' signal received, please type 'quit' to exit", <-shutdownSignalChan)
		}
	}(ctx)
	return ctx
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		newOrg := oy + dy
		if err := v.SetOrigin(ox, newOrg); err != nil {
			//if dy > 0 {
			//	v.Autoscroll = true
			//}
			return err
		}
	}
	return nil
}

func scrollToEnd(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = true
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

var cmdStack []string
var curNdx = -1

func switchView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "stdin" {
		m, err := g.SetCurrentView("strategies")
		m.Frame = true
		if v != nil {
			v.Frame = false
		}
		return err
	}
	s, err := g.SetCurrentView("stdin")
	LogWin.Autoscroll = true
	LogWin.Frame = false
	s.Frame = true
	return err
}

func InitKeybindings(g *gocui.Gui, cmdhandler func(string)) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			dst := LogWin
			fmt.Fprintln(dst, "Use 'quit' command for exiting application..")
			return nil
		}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("stdin", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			str := v.Buffer()
			str = strings.Trim(str, "\n")
			if len(str) > 0 {
				if str == "quit" {
					return gocui.ErrQuit
				} else {
					if len(cmdStack) == 0 || cmdStack[len(cmdStack)-1] != str {
						cmdStack = append(cmdStack, str)
					}
					curNdx = len(cmdStack)
					go cmdhandler(str)
					GlobalGui.Update(func(gui *gocui.Gui) error {
						return nil
					})

				}
			}
			v.Clear()
			v.SetCursor(0, 0)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if curNdx <= 0 {
				return nil
			}
			curNdx--
			s := cmdStack[curNdx]
			v.Clear()
			v.SetCursor(0, 0)
			fmt.Fprint(v, s)
			v.SetCursor(len(s), 0)

			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {

			if curNdx >= len(cmdStack) {
				v.Clear()
				v.SetCursor(0, 0)
				return nil
			}
			curNdx++
			if curNdx >= len(cmdStack) {
				v.Clear()
				v.SetCursor(0, 0)
				return nil
			}
			s := cmdStack[curNdx]
			v.Clear()
			v.SetCursor(0, 0)
			fmt.Fprint(v, s)
			v.SetCursor(len(s), 0)

			return nil
		}); err != nil {
		return err
	}

	_, maxY := g.Size()
	if err := g.SetKeybinding("stdin", gocui.KeyTab, gocui.ModNone, switchView); err != nil {
		return err
	}
	if err := g.SetKeybinding("strategies", gocui.KeyTab, gocui.ModNone, switchView); err != nil {
		return err
	}

	if err := g.SetKeybinding("stdin", gocui.KeyEnd, gocui.ModAlt,
		func(g *gocui.Gui, v *gocui.View) error {
			LogWin.Autoscroll = true
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyHome, gocui.ModAlt,
		func(g *gocui.Gui, v *gocui.View) error {
			LogWin.Autoscroll = false
			if err := LogWin.SetOrigin(0, 0); err != nil {
				LogWin.Autoscroll = true
				return err
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("strategies", gocui.KeyEnd, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			LogWin.Autoscroll = true
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("strategies", gocui.KeyHome, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			LogWin.Autoscroll = false
			if err := LogWin.SetOrigin(0, 0); err != nil {
				LogWin.Autoscroll = true
				return err
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyPgdn, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(LogWin, maxY-3)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyPgup, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(LogWin, -maxY+3)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("strategies", gocui.KeyPgdn, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(LogWin, maxY-3)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("strategies", gocui.KeyPgup, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(LogWin, -maxY+3)
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("stdin", gocui.KeyArrowUp, gocui.ModAlt,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(LogWin, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("stdin", gocui.KeyArrowDown, gocui.ModAlt,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(LogWin, 1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("strategies", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(LogWin, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("strategies", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(LogWin, 1)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("strategies", 0, 0, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
		v.Frame = false
		v.Wrap = true
	}

	if v, err := g.SetView("stdin", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if _, err := g.SetCurrentView("stdin"); err != nil {
			return err
		}
		g.Cursor = true
		v.Editable = true
		v.Wrap = true
	}
	return nil
}

func printPrompt(prompt string) {

	fmt.Printf("\r%s", prompt)
}

func cleanPrompt(l int) {
	fm := fmt.Sprintf("\r-%%%d.%ds", l, l)
	fmt.Printf(fm, " ")
}

// AskForConfirmation ask user for y or n
// -- source: https://gist.github.com/albrow/5882501
// utils.AskForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling utils.AskForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
// default value defines empty (just Enter) response :
// -1 "no"
// 1 "yes"
// 0 NO default response allowed
func AskForConfirmation(prompt string, dflt int) bool {

	var response string
	var err error

	if !UseGuiQ() {
		printPrompt(prompt)
		_, err = fmt.Scanln(&response)
		if err != nil {
			switch dflt {
			case 1:
				return true
			case -1:
				return false
			case 0:
				printPrompt("try again: ")
				return AskForConfirmation(prompt, dflt)
			}
		}

	} else {
		var pResponse *string
		var gotResult AtomBool

		GlobalGui.Update(func(gui *gocui.Gui) error {
			ShowConfView(GlobalGui, prompt, func(r *string) {
				pResponse = r
				gotResult.Set()
			})
			return nil
		})
		for !gotResult.Get() {
			time.Sleep(50 * time.Millisecond)
		}

		if pResponse == nil {
			return AskForConfirmation(prompt, dflt)
		}

		response = *pResponse
		if len(response) == 0 {

			switch dflt {
			case 1:
				return true
			case -1:
				return false
			case 0:
				return AskForConfirmation(prompt, dflt)
			}
		}
	}
	yops := []string{"y", "Y", "yes", "Yes", "YES"}
	nops := []string{"n", "N", "no", "No", "NO"}
	if containsString(yops, response) {
		return true
	} else if containsString(nops, response) {
		return false
	} else {
		if !UseGuiQ() {
			printPrompt("try again: ")
		}
		return AskForConfirmation(prompt, dflt)
	}
}
