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
	wnd := ui.NewWindowControl(parent,
		ui.WindowControlOpts().
			Position(pos).
			Size(sz).
			HCursor(win.HINSTANCE(0).LoadCursor(co.IDC_HAND)).
			WndExStyles(co.WS_EX_CLIENTEDGE),
	)

	me := &Tracker{
		wnd: wnd,
	}

	me.events()
	return me
}

func (me *Tracker) events() {
	me.wnd.On().WmPaint(func() {
		hwnd := me.wnd.Hwnd()

		ps := win.PAINTSTRUCT{}
		hdc := hwnd.BeginPaint(&ps)
		defer hwnd.EndPaint(&ps)

		fillColor := win.GetSysColor(co.COLOR_ACTIVECAPTION)

		myPen := win.CreatePen(co.PS_SOLID, 1, fillColor)
		defer myPen.DeleteObject()
		defPen := hdc.SelectObjectPen(myPen)
		defer hdc.SelectObjectPen(defPen)

		myBrush := win.CreateSolidBrush(fillColor)
		defer myBrush.DeleteObject()
		defBrush := hdc.SelectObjectBrush(myBrush)
		defer hdc.SelectObjectBrush(defBrush)

		rcClient := hwnd.GetClientRect()
		hdc.Rectangle(0, 0, int32(float32(rcClient.Right)*me.elapsed), rcClient.Bottom)
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
