package main

import (
	"runtime"

	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"github.com/rodrigocfd/windigo/win/com/shell"
	"github.com/rodrigocfd/windigo/win/com/shell/shellco"
)

func main() {
	runtime.LockOSThread()

	win.CoInitializeEx(co.COINIT_APARTMENTTHREADED)
	defer win.CoUninitialize()

	m := NewMain()
	m.Run()
}

const CMD_OPEN int = 20_000

type Main struct {
	wnd     ui.WindowMain
	pic     *Picture
	tracker *Tracker
	resz    ui.Resizer
}

func NewMain() *Main {
	wnd := ui.NewWindowMain(
		ui.WindowMainOpts().
			Title("The playback").
			IconId(101).
			AccelTable(ui.NewAcceleratorTable().
				AddChar('O', co.ACCELF_CONTROL, CMD_OPEN)).
			ClientArea(win.SIZE{Cx: 500, Cy: 300}).
			WndStyles(co.WS_CAPTION | co.WS_SYSMENU | co.WS_CLIPCHILDREN |
				co.WS_BORDER | co.WS_VISIBLE | co.WS_MINIMIZEBOX |
				co.WS_MAXIMIZEBOX | co.WS_SIZEBOX),
	)

	me := &Main{
		wnd:     wnd,
		pic:     NewPicture(wnd, win.POINT{X: 10, Y: 10}, win.SIZE{Cx: 480, Cy: 250}),
		tracker: NewTracker(wnd, win.POINT{X: 10, Y: 270}, win.SIZE{Cx: 480, Cy: 20}),
		resz:    ui.NewResizer(wnd),
	}

	me.resz.Add(ui.RESZ_RESIZE, ui.RESZ_RESIZE, me.pic.wnd).
		Add(ui.RESZ_RESIZE, ui.RESZ_REPOS, me.tracker.wnd)
	me.events()
	return me
}

func (me *Main) Run() {
	defer me.pic.Free()

	me.wnd.RunAsMain()
}

func (me *Main) events() {
	me.wnd.On().WmCreate(func(p wm.Create) int {
		me.wnd.Hwnd().SetTimer(1, 100, func(msElapsed uint32) {
			// memStats := runtime.MemStats{}
			// runtime.ReadMemStats(&memStats)

			// me.wnd.Hwnd().SetWindowText(
			// 	fmt.Sprintf("%s / Alloc: %s, cycles: %d, next: %s",
			// 		me.pic.CurrentPosDurFmt(),
			// 		win.Str.FmtBytes(memStats.HeapAlloc),
			// 		memStats.NumGC,
			// 		win.Str.FmtBytes(memStats.NextGC)))

			me.wnd.Hwnd().SetWindowText(me.pic.CurrentPosDurFmt())

			me.tracker.SetElapsed(float32(me.pic.CurrentPos()) / float32(me.pic.Duration()))
		})
		return 0
	})

	me.wnd.On().WmCommandAccelMenu(CMD_OPEN, func(_ wm.Command) {
		me.pic.Pause()

		fod := shell.CoCreateIFileOpenDialog(co.CLSCTX_INPROC_SERVER)
		defer fod.Release()

		flags := fod.GetOptions()
		fod.SetOptions(flags | shellco.FOS_FORCEFILESYSTEM | shellco.FOS_FILEMUSTEXIST)

		fod.SetFileTypes([]shell.FilterSpec{
			{Name: "All video files", Spec: "*.mkv;*.mp4"},
			{Name: "Matroska", Spec: "*.mkv"},
			{Name: "MPEG-4", Spec: "*.mp4"},
			{Name: "Anything", Spec: "*.*"},
		})
		fod.SetFileTypeIndex(0)

		if fod.Show(me.wnd.Hwnd()) {
			shi := fod.GetResult()
			defer shi.Release()

			me.pic.StartPlayback(shi.GetDisplayName(shellco.SIGDN_FILESYSPATH))
		}
	})

	me.wnd.On().WmCommandAccelMenu(int(co.ID_CANCEL), func(_ wm.Command) {
		me.wnd.Hwnd().SendMessage(co.WM_CLOSE, 0, 0)
	})

	me.tracker.OnClick(func(pct float32) {
		me.pic.SetCurrentPos(int(float32(me.pic.Duration()) * pct))
	})
}
