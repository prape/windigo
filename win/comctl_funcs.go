//go:build windows

package win

import (
	"syscall"
	"unsafe"

	"github.com/rodrigocfd/windigo/internal/dll"
	"github.com/rodrigocfd/windigo/internal/utl"
	"github.com/rodrigocfd/windigo/win/co"
	"github.com/rodrigocfd/windigo/win/wstr"
)

// [ImageList_DragMove] function.
//
// [ImageList_DragMove]: https://learn.microsoft.com/en-us/windows/win32/api/commctrl/nf-commctrl-imagelist_dragmove
func ImageListDragMove(x, y int) error {
	ret, _, _ := syscall.SyscallN(
		dll.Load(dll.COMCTL32, &_ImageList_DragMove, "ImageList_DragMove"),
		uintptr(int32(x)),
		uintptr(int32(y)))
	return utl.ZeroAsSysInvalidParm(ret)
}

var _ImageList_DragMove *syscall.Proc

// [ImageList_DragShowNolock] function.
//
// [ImageList_DragShowNolock]: https://learn.microsoft.com/en-us/windows/win32/api/commctrl/nf-commctrl-imagelist_dragshownolock
func ImageListDragShowNolock(show bool) error {
	ret, _, _ := syscall.SyscallN(
		dll.Load(dll.COMCTL32, &_ImageList_DragShowNolock, "ImageList_DragShowNolock"),
		utl.BoolToUintptr(show))
	return utl.ZeroAsSysInvalidParm(ret)
}

var _ImageList_DragShowNolock *syscall.Proc

// [ImageList_DrawIndirect] function.
//
// [ImageList_DrawIndirect]: https://learn.microsoft.com/en-us/windows/win32/api/commctrl/nf-commctrl-imagelist_drawindirect
func ImageListDrawIndirect(imldp *IMAGELISTDRAWPARAMS) error {
	ret, _, _ := syscall.SyscallN(
		dll.Load(dll.COMCTL32, &_ImageList_DrawIndirect, "ImageList_DrawIndirect"),
		uintptr(unsafe.Pointer(imldp)))
	return utl.ZeroAsSysInvalidParm(ret)
}

var _ImageList_DrawIndirect *syscall.Proc

// [ImageList_EndDrag] function.
//
// [ImageList_EndDrag]: https://learn.microsoft.com/en-us/windows/win32/api/commctrl/nf-commctrl-imagelist_enddrag
func ImageListEndDrag() {
	syscall.SyscallN(
		dll.Load(dll.COMCTL32, &_ImageList_EndDrag, "ImageList_EndDrag"))
}

var _ImageList_EndDrag *syscall.Proc

// [InitCommonControls] function.
//
// [InitCommonControls]: https://learn.microsoft.com/en-us/windows/win32/api/commctrl/nf-commctrl-initcommoncontrols
func InitCommonControls() {
	syscall.SyscallN(
		dll.Load(dll.COMCTL32, &_InitCommonControls, "InitCommonControls"))
}

var _InitCommonControls *syscall.Proc

// [InitCommonControlsEx] function.
//
// [InitCommonControlsEx]: https://learn.microsoft.com/en-us/windows/win32/api/commctrl/nf-commctrl-initcommoncontrolsex
func InitCommonControlsEx(icc co.ICC) error {
	var iccx _INITCOMMONCONTROLSEX
	iccx.SetDwSize()
	iccx.DwICC = icc

	ret, _, _ := syscall.SyscallN(
		dll.Load(dll.COMCTL32, &_InitCommonControlsEx, "InitCommonControlsEx"),
		uintptr(unsafe.Pointer(&iccx)))
	return utl.ZeroAsSysInvalidParm(ret)
}

var _InitCommonControlsEx *syscall.Proc

// [InitMUILanguage] function.
//
// [InitMUILanguage]: https://learn.microsoft.com/en-us/windows/win32/api/commctrl/nf-commctrl-initmuilanguage
func InitMUILanguage(lang LANGID) {
	syscall.SyscallN(
		dll.Load(dll.COMCTL32, &_InitMUILanguage, "InitMUILanguage"),
		uintptr(lang))
}

var _InitMUILanguage *syscall.Proc

// [TaskDialogIndirect] function.
//
// Example:
//
//	var hWnd win.HWND // initialized somewhere
//
//	_, _ = win.TaskDialogIndirect(win.TASKDIALOGCONFIG{
//		HwndParent:      hWnd,
//		WindowTitle:     "Title",
//		MainInstruction: "Caption",
//		Content:         "Body",
//		HMainIcon:       win.TdcIconTdi(co.TDICON_INFORMATION),
//		CommonButtons:   co.TDCBF_OK,
//		Flags: co.TDF_ALLOW_DIALOG_CANCELLATION |
//			co.TDF_POSITION_RELATIVE_TO_WINDOW,
//	})
//
// [TaskDialogIndirect]: https://learn.microsoft.com/en-us/windows/win32/api/commctrl/nf-commctrl-taskdialogindirect
func TaskDialogIndirect(taskConfig TASKDIALOGCONFIG) (co.ID, error) {
	wbuf := wstr.NewBufEncoder() // to keep all strings used in the call
	defer wbuf.Free()

	tdcBuf := NewVecSized(160, byte(0)) // packed TASKDIALOGCONFIG is 160 bytes
	defer tdcBuf.Free()

	btnsBuf := NewVec[[12]byte]() // packed TASKDIALOG_BUTTON is 12 bytes, we can have many
	defer btnsBuf.Free()

	taskConfig.serialize(&wbuf, &tdcBuf, &btnsBuf)

	pPnButtons := NewVecSized(3, int32(0)) // button, radio and check values returned
	defer pPnButtons.Free()

	ret, _, _ := syscall.SyscallN(
		dll.Load(dll.COMCTL32, &_TaskDialogIndirect, "TaskDialogIndirect"),
		uintptr(tdcBuf.UnsafePtr()),
		uintptr(unsafe.Pointer(pPnButtons.Get(0))),
		uintptr(unsafe.Pointer(pPnButtons.Get(1))),
		uintptr(unsafe.Pointer(pPnButtons.Get(2))))
	if hr := co.HRESULT(ret); hr != co.HRESULT_S_OK {
		return co.ID(0), hr
	}

	return co.ID(*pPnButtons.Get(0)), nil
}

var _TaskDialogIndirect *syscall.Proc
