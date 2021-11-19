package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"github.com/rodrigocfd/windigo/win/com/shell"
	"github.com/rodrigocfd/windigo/win/com/shell/shellco"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			win.HWND(0).MessageBox(
				fmt.Sprintf("PANIC @ %v\n\n%v\n\n%s",
					time.Now(), r, string(debug.Stack())),
				"Panic", co.MB_ICONERROR)
		}
	}()
	runtime.LockOSThread()

	win.CoInitializeEx(co.COINIT_APARTMENTTHREADED)
	defer win.CoUninitialize()

	m := NewMain()
	m.Run()

	// clsId, _ := win.CLSIDFromProgID("Excel.Application")
	// root := win.CoCreateInstance(clsId, nil, co.CLSCTX_SERVER, co.IID_IUNKNOWN)
	// defer root.Release()
	// excel := autom.NewIDispatch(root.QueryInterface(automco.IID_IDispatch))
	// defer excel.Release()
	// for _, f := range excel.ListFunctions() {
	// 	println(f.Name, f.NumParams, f.FuncKind, f.Flags)
	// }
}

const (
	CMD_OPEN = iota + 20_000
	CMD_ABOUT
)

// Main application window.
type Main struct {
	wnd     ui.WindowMain
	pic     *Picture
	tracker *Tracker
}

func NewMain() *Main {
	wnd := ui.NewWindowMain(
		ui.WindowMainOpts().
			Title("The playback").
			IconId(101).
			AccelTable(ui.NewAcceleratorTable().
				AddChar('O', co.ACCELF_CONTROL, CMD_OPEN).
				AddKey(co.VK_F1, co.ACCELF_NONE, CMD_ABOUT)).
			ClientArea(win.SIZE{Cx: 600, Cy: 270}).
			WndStyles(co.WS_CAPTION | co.WS_SYSMENU | co.WS_CLIPCHILDREN |
				co.WS_BORDER | co.WS_VISIBLE | co.WS_MINIMIZEBOX |
				co.WS_MAXIMIZEBOX | co.WS_SIZEBOX).
			WndExStyles(co.WS_EX_ACCEPTFILES),
	)

	me := &Main{
		wnd: wnd,
		pic: NewPicture(wnd,
			win.POINT{X: 2, Y: 2},
			win.SIZE{Cx: 596, Cy: 250},
			ui.HORZ_RESIZE, ui.VERT_RESIZE),
		tracker: NewTracker(wnd,
			win.POINT{X: 2, Y: 252},
			win.SIZE{Cx: 596, Cy: 16},
			ui.HORZ_RESIZE, ui.VERT_REPOS),
	}

	me.events()
	return me
}

func (me *Main) Run() {
	defer me.pic.Free()

	me.wnd.RunAsMain()
}

func (me *Main) events() {
	me.wnd.On().WmCreate(func(p wm.Create) int {
		me.wnd.Hwnd().SetTimerCallback(100, func(_ uintptr) {
			me.wnd.Hwnd().SetWindowText(me.pic.CurrentPosDurFmt())
			me.tracker.SetElapsed(float32(me.pic.CurrentPos()) / float32(me.pic.Duration()))
		})
		return 0
	})

	me.wnd.On().WmDropFiles(func(p wm.DropFiles) {
		droppedFiles := p.Hdrop().ListFilesAndFinish()
		if win.Path.HasExtension(droppedFiles[0], ".avi", ".mkv", ".mp4") {
			me.pic.StartPlayback(droppedFiles[0])
		}
	})

	me.wnd.On().WmCommandAccelMenu(CMD_OPEN, func(_ wm.Command) {
		me.pic.Pause()

		fod := shell.NewIFileOpenDialog(
			win.CoCreateInstance(
				shellco.CLSID_FileOpenDialog, nil,
				co.CLSCTX_INPROC_SERVER,
				shellco.IID_IFileOpenDialog),
		)
		defer fod.Release()

		fod.SetOptions(fod.GetOptions() |
			shellco.FOS_FORCEFILESYSTEM | shellco.FOS_FILEMUSTEXIST)

		fod.SetFileTypes([]shell.FilterSpec{
			{Name: "All video files", Spec: "*.avi;*.mkv;*.mp4"},
			{Name: "AVI", Spec: "*.avi"},
			{Name: "Matroska", Spec: "*.mkv"},
			{Name: "MPEG-4", Spec: "*.mp4"},
			{Name: "Anything", Spec: "*.*"},
		})
		fod.SetFileTypeIndex(1)

		// shiDir, _ := shell.NewShellItem(win.GetCurrentDirectory())
		// defer shiDir.Release()
		// fod.SetFolder(&shiDir)

		if fod.Show(me.wnd.Hwnd()) {
			me.pic.StartPlayback(fod.GetResultDisplayName(shellco.SIGDN_FILESYSPATH))
		}
	})

	me.wnd.On().WmCommandAccelMenu(CMD_ABOUT, func(_ wm.Command) {
		me.pic.Pause()

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		tdc := win.TASKDIALOGCONFIG{
			HwndParent:         me.wnd.Hwnd(),
			DwFlags:            co.TDF_ALLOW_DIALOG_CANCELLATION,
			DwCommonButtons:    co.TDCBF_OK,
			PszWindowTitle:     "About",
			HMainIcon:          win.TdcIconTdi(co.TD_ICON_INFORMATION),
			PszMainInstruction: "Playback",
			PszContent: fmt.Sprintf(
				"Windigo experimental playback application.\n\n"+
					"Objects mem: %s\n"+
					"Reserved sys: %s\n"+
					"Idle spans: %s\n"+
					"GC cycles: %d\n"+
					"Next GC: %s",
				win.Str.FmtBytes(memStats.HeapAlloc),
				win.Str.FmtBytes(memStats.HeapSys),
				win.Str.FmtBytes(memStats.HeapIdle),
				memStats.NumGC,
				win.Str.FmtBytes(memStats.NextGC)),
		}
		win.TaskDialogIndirect(&tdc)
	})

	me.wnd.On().WmCommandAccelMenu(int(co.ID_CANCEL), func(_ wm.Command) { // close on ESC
		me.wnd.Hwnd().SendMessage(co.WM_CLOSE, 0, 0)
	})

	me.tracker.OnClick(func(pct float32) {
		me.pic.SetCurrentPos(int(float32(me.pic.Duration()) * pct))
	})

	me.tracker.OnSpace(func() {
		me.pic.TogglePlayPause()
	})

	me.tracker.OnLeftRight(func(key co.VK) {
		if key == co.VK_LEFT {
			me.pic.BackwardSecs(10)
		} else if key == co.VK_RIGHT {
			me.pic.ForwardSecs(10)
		}
	})
}
