package main

import (
	"fmt"
	"runtime"

	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"github.com/rodrigocfd/windigo/win/com/shell"
)

func main() {
	win.CoInitializeEx(co.COINIT_APARTMENTTHREADED)
	defer win.CoUninitialize()

	m := NewMain()
	m.Run()
}

const CMD_OPEN int = 20_000

type Main struct {
	wnd ui.WindowMain
	pic *Picture
}

func NewMain() *Main {
	wnd := ui.NewWindowMainOpts(ui.WindowMainOpts{
		Title:  "The playback",
		IconId: 101,
		AccelTable: ui.NewAcceleratorTable().
			AddChar('O', co.ACCELF_CONTROL, CMD_OPEN),
		ClientAreaSize: win.SIZE{Cx: 500, Cy: 300},
	})

	me := Main{
		wnd: wnd,
		pic: NewPicture(wnd, win.POINT{X: 10, Y: 10}, win.SIZE{Cx: 480, Cy: 280}),
	}

	me.events()
	return &me
}

func (me *Main) Run() {
	defer me.pic.Free()

	me.wnd.RunAsMain()
}

var memStats runtime.MemStats // cache

func (me *Main) events() {
	me.wnd.On().WmCreate(func(p wm.Create) int {
		me.wnd.Hwnd().SetTimer(1, 500, func(msElapsed uint32) {
			runtime.ReadMemStats(&memStats)
			me.wnd.Hwnd().SetWindowText(
				fmt.Sprintf("%s / Alloc: %s, cycles: %d, next: %s",
					me.pic.CurrentTime(),
					win.Str.FmtBytes(memStats.HeapAlloc),
					memStats.NumGC,
					win.Str.FmtBytes(memStats.NextGC)))
		})
		return 0
	})

	me.wnd.On().WmCommandAccelMenu(CMD_OPEN, func(_ wm.Command) {
		vidPath, ok := ui.Prompt.OpenSingleFile(me.wnd, []shell.FilterSpec{
			{Name: "All video files", Spec: "*.mkv;*.mp4"},
			{Name: "Matroska", Spec: "*.mkv"},
			{Name: "MPEG-4", Spec: "*.mp4"},
			{Name: "Anything", Spec: "*.*"},
		})
		if ok {
			me.pic.StartPlayback(vidPath)
		}
	})

	me.wnd.On().WmCommandAccelMenu(int(co.ID_CANCEL), func(_ wm.Command) {
		me.wnd.Hwnd().SendMessage(co.WM_CLOSE, 0, 0)
	})
}
