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

func NewPicture(parent ui.AnyParent, pos win.POINT, sz win.SIZE) *Picture {
	wnd := ui.NewWindowControl(parent,
		ui.WindowControlOpts().
			Position(pos).
			Size(sz),
	)

	me := &Picture{
		wnd: wnd,
	}

	me.events()
	return me
}

func (me *Picture) Free() {
	if me.mediaCtrl.Ppv != nil {
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
		ps := win.PAINTSTRUCT{}
		me.wnd.Hwnd().BeginPaint(&ps)
		defer me.wnd.Hwnd().EndPaint(&ps)

		if me.controllerEvr.Ppv != nil {
			me.controllerEvr.RepaintVideo()
		}
	})

	me.wnd.On().WmSize(func(p wm.Size) {
		if me.controllerEvr.Ppv != nil {
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

	me.graphBuilder = dshow.NewIGraphBuilder(co.CLSCTX_INPROC_SERVER)
	me.vmr = dshow.NewEnhancedVideoRenderer(co.CLSCTX_INPROC_SERVER)
	me.graphBuilder.AddFilter(&me.vmr, "EVR")

	getSvc := me.vmr.QueryIMFGetService()
	defer getSvc.Release()

	me.controllerEvr = getSvc.GetIMFVideoDisplayControl()
	if e := me.controllerEvr.SetVideoWindow(me.wnd.Hwnd()); e != nil {
		panic(e)
	}
	if e := me.controllerEvr.SetAspectRatioMode(dshowco.MFVideoARMode_PreservePicture); e != nil {
		panic(e)
	}

	me.mediaCtrl = me.graphBuilder.QueryIMediaControl()
	me.mediaSeek = me.graphBuilder.QueryIMediaSeeking()
	me.basicAudio = me.graphBuilder.QueryIBasicAudio()

	if e := me.graphBuilder.RenderFile(vidPath); e != nil {
		panic(e)
	}

	rc := me.wnd.Hwnd().GetWindowRect()
	me.wnd.Hwnd().ScreenToClientRc(&rc)
	me.controllerEvr.SetVideoPosition(nil, &rc)

	me.mediaCtrl.Run()
}

func (me *Picture) Pause() {
	if me.mediaCtrl.Ppv != nil {
		state, _ := me.mediaCtrl.GetState(-1)
		if state == dshowco.FILTER_STATE_State_Running {
			me.mediaCtrl.Pause()
		}
	}
}

func (me *Picture) TogglePlayPause() {
	if me.mediaCtrl.Ppv != nil {
		state, _ := me.mediaCtrl.GetState(-1)
		if state == dshowco.FILTER_STATE_State_Running {
			me.mediaCtrl.Pause()
		} else {
			me.mediaCtrl.Run()
		}
	}
}

func (me *Picture) Duration() (secs int) {
	if me.mediaSeek.Ppv == nil {
		secs = 0
	} else {
		secs = int(me.mediaSeek.GetDuration() / time.Second)
	}
	return
}

func (me *Picture) SetCurrentPos(secs int) {
	if me.mediaSeek.Ppv != nil {
		me.mediaSeek.SetPositions(
			time.Duration(secs)*time.Second, dshowco.SEEKING_FLAGS_AbsolutePositioning,
			0, dshowco.SEEKING_FLAGS_NoPositioning)
	}
}

func (me *Picture) CurrentPos() (secs int) {
	if me.mediaSeek.Ppv == nil {
		secs = 0
	} else {
		secs = int(me.mediaSeek.GetCurrentPosition() / time.Second)
	}
	return
}

func (me *Picture) CurrentPosDurFmt() string {
	if me.mediaSeek.Ppv == nil {
		return "NO VIDEO"
	}

	stCurPos := win.SYSTEMTIME{}
	stCurPos.FromDuration(me.mediaSeek.GetCurrentPosition())

	stDur := win.SYSTEMTIME{}
	stDur.FromDuration(me.mediaSeek.GetDuration())

	return fmt.Sprintf("%d:%02d:%02d of %d:%02d:%02d",
		stCurPos.WHour, stCurPos.WMinute, stCurPos.WSecond,
		stDur.WHour, stDur.WMinute, stDur.WSecond)
}

func (me *Picture) ForwardSecs(secs int) {
	if me.mediaSeek.Ppv != nil {
		newSecs := me.CurrentPos() + secs
		duration := me.Duration()
		if newSecs >= duration {
			newSecs = duration - 1 // max pos
		}
		me.SetCurrentPos(newSecs)
	}
}

func (me *Picture) BackwardSecs(secs int) {
	if me.mediaSeek.Ppv != nil {
		newSecs := me.CurrentPos() - secs
		if newSecs < 0 {
			newSecs = 0 // min pos
		}
		me.SetCurrentPos(newSecs)
	}
}
