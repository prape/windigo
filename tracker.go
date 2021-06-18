package main

import (
	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/ui/wm"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
)

type Tracker struct {
	wnd      ui.WindowControl
	funClick func(pct float32)
	elapsed  float32
}

func NewTracker(parent ui.AnyParent, pos win.POINT, sz win.SIZE) *Tracker {
	wnd := ui.NewWindowControlRaw(parent, ui.WindowControlRawOpts{
		Position: pos,
		Size:     sz,
		HCursor:  win.HINSTANCE(0).LoadCursor(co.IDC_HAND),
		ExStyles: co.WS_EX(0),
	})

	me := &Tracker{
		wnd: wnd,
	}

	me.events()
	return me
}

func (me *Tracker) events() {
	me.wnd.On().WmPaint(func() {
		ps := win.PAINTSTRUCT{}
		hdc := me.wnd.Hwnd().BeginPaint(&ps)
		defer me.wnd.Hwnd().EndPaint(&ps)

		myBrush := win.CreateSolidBrush(win.GetSysColor(co.COLOR_ACTIVECAPTION))
		defer myBrush.DeleteObject()

		oldBrush := hdc.SelectObjectBrush(myBrush)
		defer hdc.SelectObjectBrush(oldBrush)

		rcClient := me.wnd.Hwnd().GetClientRect()
		hdc.Rectangle(0, 0, int32(float32(rcClient.Right)*me.elapsed)-1, rcClient.Bottom-1)
	})

	me.wnd.On().WmLButtonDown(func(p wm.Mouse) {
		if me.funClick != nil {
			rcClient := me.wnd.Hwnd().GetClientRect()
			pct := float32(p.Pos().X) / float32(rcClient.Right)
			me.funClick(pct)
		}
	})
}

func (me *Tracker) OnClick(fun func(pct float32)) {
	me.funClick = fun
}

func (me *Tracker) SetElapsed(pct float32) {
	me.elapsed = pct
	me.wnd.Hwnd().InvalidateRect(nil, true)
}
