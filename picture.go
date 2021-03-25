package main

import (
	"fmt"

	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"github.com/rodrigocfd/windigo/win/com/dshow"
)

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
	wnd := ui.NewWindowControlOpts(parent, ui.WindowControlOpts{
		Position: pos,
		Size:     sz,
	})

	me := Picture{
		wnd: wnd,
	}

	me.events()
	return &me
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

	me.wnd.On().WmLButtonDown(func(_ wm.Mouse) {
		state, _ := me.mediaCtrl.GetState(-1)
		if state == co.FILTER_STATE_State_Running {
			me.mediaCtrl.Pause()
		} else {
			me.mediaCtrl.Run()
		}
	})
}

func (me *Picture) CurrentTime() string {
	if me.mediaSeek.Ppv == nil {
		return ""
	}

	curPos, _ := me.mediaSeek.GetPositions()
	curPos /= 1000 * 100 // millisecond * nanosecond
	return fmt.Sprintf("%d:%d:%d",
		(curPos-(curPos%(60*60)))/60*60,
		curPos-(curPos%60),
		curPos)
}

func (me *Picture) StartPlayback(vidPath string) {
	me.Free()

	me.graphBuilder = dshow.CoCreateIGraphBuilder(co.CLSCTX_INPROC_SERVER)
	me.vmr = dshow.CoCreateEnhancedVideoRenderer(co.CLSCTX_INPROC_SERVER)
	if e := me.graphBuilder.AddFilter(&me.vmr, "EVR"); e != nil {
		panic(e)
	}

	getSvc := me.vmr.QueryIMFGetService()
	defer getSvc.Release()

	me.controllerEvr = getSvc.GetServiceIMFVideoDisplayControl()
	if e := me.controllerEvr.SetVideoWindow(me.wnd.Hwnd()); e != nil {
		panic(e)
	}
	if e := me.controllerEvr.SetAspectRatioMode(co.MFVideoARMode_PreservePicture); e != nil {
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
