package main

import (
	"fmt"
	"time"

	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"github.com/rodrigocfd/windigo/win/com/dshow"
	"github.com/rodrigocfd/windigo/win/com/dshow/dshowco"
)

// Child window which renders the video.
type Picture struct {
	wnd ui.WindowControl

	graphBuilder  dshow.IGraphBuilder
	vmr           dshow.IBaseFilter
	controllerEvr dshow.IMFVideoDisplayControl
	mediaCtrl     dshow.IMediaControl
	mediaSeek     dshow.IMediaSeeking
	basicAudio    dshow.IBasicAudio
}

func NewPicture(
	parent ui.AnyParent, pos win.POINT, sz win.SIZE,
	horz ui.HORZ, vert ui.VERT) *Picture {

	wnd := ui.NewWindowControl(parent,
		ui.WindowControlOpts().
			WndExStyles(co.WS_EX_NONE).
			Position(pos).
			Size(sz).
			Horz(horz).
			Vert(vert),
	)

	me := &Picture{
		wnd: wnd,
	}

	me.events()
	return me
}

func (me *Picture) Free() {
	if me.mediaCtrl.Ptr() != nil {
		me.mediaCtrl.Stop()
	}

	me.basicAudio.Release()
	me.mediaSeek.Release()
	me.mediaCtrl.Release()
	me.controllerEvr.Release()
	me.vmr.Release()
	me.graphBuilder.Release()
}

func (me *Picture) events() {
	me.wnd.On().WmPaint(func() {
		var ps win.PAINTSTRUCT
		me.wnd.Hwnd().BeginPaint(&ps)
		defer me.wnd.Hwnd().EndPaint(&ps)

		if me.controllerEvr.Ptr() != nil {
			me.controllerEvr.RepaintVideo()
		}
	})

	me.wnd.On().WmSize(func(p wm.Size) {
		if me.controllerEvr.Ptr() != nil {
			rc := me.wnd.Hwnd().GetWindowRect()
			me.wnd.Hwnd().ScreenToClientRc(&rc)
			me.controllerEvr.SetVideoPosition(nil, &rc)
		}
	})

	me.wnd.On().WmLButtonDown(func(_ wm.Mouse) {
		me.wnd.Hwnd().SetFocus()
		me.TogglePlayPause()
	})

	me.wnd.On().WmKeyDown(func(p wm.Key) {
		if p.VirtualKeyCode() == co.VK_SPACE {
			me.TogglePlayPause()
		}
	})

	me.wnd.On().WmGetDlgCode(func(p wm.GetDlgCode) co.DLGC {
		if p.VirtualKeyCode() == co.VK_LEFT {
			me.BackwardSecs(10)
			return co.DLGC_WANTARROWS
		} else if p.VirtualKeyCode() == co.VK_RIGHT {
			me.ForwardSecs(10)
			return co.DLGC_WANTARROWS
		}
		return co.DLGC_NONE
	})
}

func (me *Picture) StartPlayback(vidPath string) {
	me.Free()

	me.graphBuilder = dshow.NewIGraphBuilder(
		win.CoCreateInstance(
			dshowco.CLSID_FilterGraph, nil,
			co.CLSCTX_INPROC_SERVER,
			dshowco.IID_IGraphBuilder),
	)
	me.vmr = dshow.NewIBaseFilter(
		win.CoCreateInstance(
			dshowco.CLSID_EnhancedVideoRenderer, nil,
			co.CLSCTX_INPROC_SERVER,
			dshowco.IID_IBaseFilter),
	)
	me.graphBuilder.AddFilter(&me.vmr, "EVR")

	getSvc := dshow.NewIMFGetService(
		me.vmr.QueryInterface(dshowco.IID_IMFGetService),
	)
	defer getSvc.Release()

	me.controllerEvr = dshow.NewIMFVideoDisplayControl(
		getSvc.GetService(
			win.GuidFromClsid(dshowco.CLSID_MR_VideoRenderService),
			win.GuidFromIid(dshowco.IID_IMFVideoDisplayControl),
		),
	)

	if e := me.controllerEvr.SetVideoWindow(me.wnd.Hwnd()); e != nil {
		panic(e)
	}
	if e := me.controllerEvr.SetAspectRatioMode(dshowco.MFVideoARMode_PreservePicture); e != nil {
		panic(e)
	}

	me.mediaCtrl = dshow.NewIMediaControl(
		me.graphBuilder.QueryInterface(dshowco.IID_IMediaControl),
	)
	me.mediaSeek = dshow.NewIMediaSeeking(
		me.graphBuilder.QueryInterface(dshowco.IID_IMediaSeeking),
	)
	me.basicAudio = dshow.NewIBasicAudio(
		me.graphBuilder.QueryInterface(dshowco.IID_IBasicAudio),
	)

	if e := me.graphBuilder.RenderFile(vidPath); e != nil {
		panic(e)
	}

	rc := me.wnd.Hwnd().GetWindowRect()
	me.wnd.Hwnd().ScreenToClientRc(&rc)
	me.controllerEvr.SetVideoPosition(nil, &rc)

	me.mediaCtrl.Run()
}

func (me *Picture) Pause() {
	if me.mediaCtrl.Ptr() != nil {
		state, _ := me.mediaCtrl.GetState(-1)
		if state == dshowco.FILTER_STATE_State_Running {
			me.mediaCtrl.Pause()
		}
	}
}

func (me *Picture) TogglePlayPause() {
	if me.mediaCtrl.Ptr() != nil {
		state, _ := me.mediaCtrl.GetState(-1)
		if state == dshowco.FILTER_STATE_State_Running {
			me.mediaCtrl.Pause()
		} else {
			me.mediaCtrl.Run()
		}
	}
}

func (me *Picture) Duration() (secs int) {
	if me.mediaSeek.Ptr() == nil {
		return 0
	} else {
		return int(me.mediaSeek.GetDuration() / time.Second)
	}
}

func (me *Picture) SetCurrentPos(secs int) {
	if me.mediaSeek.Ptr() != nil {
		me.mediaSeek.SetPositions(
			time.Duration(secs)*time.Second, dshowco.SEEKING_FLAGS_AbsolutePositioning,
			0, dshowco.SEEKING_FLAGS_NoPositioning)
	}
}

func (me *Picture) CurrentPos() (secs int) {
	if me.mediaSeek.Ptr() == nil {
		secs = 0
	} else {
		secs = int(me.mediaSeek.GetCurrentPosition() / time.Second)
	}
	return
}

func (me *Picture) CurrentPosDurFmt() string {
	if me.mediaSeek.Ptr() == nil {
		return "NO VIDEO"

	} else {
		var stCurPos win.SYSTEMTIME
		stCurPos.FromDuration(me.mediaSeek.GetCurrentPosition())

		var stDur win.SYSTEMTIME
		stDur.FromDuration(me.mediaSeek.GetDuration())

		return fmt.Sprintf("%d:%02d:%02d of %d:%02d:%02d",
			stCurPos.WHour, stCurPos.WMinute, stCurPos.WSecond,
			stDur.WHour, stDur.WMinute, stDur.WSecond)
	}
}

func (me *Picture) ForwardSecs(secs int) {
	if me.mediaSeek.Ptr() != nil {
		newSecs := me.CurrentPos() + secs
		duration := me.Duration()
		if newSecs >= duration {
			newSecs = duration - 1 // max pos
		}
		me.SetCurrentPos(newSecs)
	}
}

func (me *Picture) BackwardSecs(secs int) {
	if me.mediaSeek.Ptr() != nil {
		newSecs := me.CurrentPos() - secs
		if newSecs < 0 {
			newSecs = 0 // min pos
		}
		me.SetCurrentPos(newSecs)
	}
}
