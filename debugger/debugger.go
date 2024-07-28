package debugger

import (
	"fmt"
	"runtime"
	"strconv"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/xSaCh/intcode/vm"
)

var (
	vmIns     *vm.IntcodeVM
	backend   imgui.Backend[imgui.GLFWWindowFlags]
	barValues []int64
)

func rawMemory() {
	memMaxWidth := len(strconv.Itoa(len(vmIns.Memory)))
	imgui.Begin("Raw Memory")
	if imgui.BeginTable("tbl_raw_mem", 2) {

		imgui.PushStyleColorU32(imgui.ColText, 0xFF5F5F5F)
		for i := 0; i < len(vmIns.Memory); i++ {
			imgui.TableNextRow()

			if i == 2 {
				imgui.PushStyleColorU32(imgui.ColText, 0xFFFFFFFF)
			}
			for j := 0; j < 2; j++ {
				imgui.TableSetColumnIndex(int32(j))
				if j == 0 {
					// imgui.PushStyleColorU32(imgui.ColText, 0xFFFFFFFF)
					imgui.Text(fmt.Sprintf("[ %0*d ]\t", memMaxWidth+1, i))
					// imgui.PopStyleColor()

				} else {
					imgui.Text(fmt.Sprintf("%d\t", vmIns.Memory[i]))
				}
			}
			if i == 2 {
				imgui.PopStyleColor()
			}
		}
		imgui.PopStyleColor()

	}
	imgui.EndTable()
	imgui.End()
}

func disassembledIntcode() {
	data, err := GetFormattedMemory(vmIns.Memory)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	imgui.Begin("Disassembled Intcode")
	if imgui.BeginTable("tbl_dis_int", 4) {

		imgui.PushStyleColorU32(imgui.ColText, 0xFF5F5F5F)
		for i := 0; i < len(data); i++ {
			imgui.TableNextRow()

			if i == 2 {
				imgui.PushStyleColorU32(imgui.ColText, 0xFFFFFFFF)
			}
			for j := 0; j < 4; j++ {
				imgui.TableSetColumnIndex(int32(j))
				// if i == 2 {
				// 	imgui.TableSetBgColor(imgui.TableBgTargetCellBg, 0xFFFF4455)
				// }
				if j < len(data[i]) {
					imgui.Text(fmt.Sprintf("%s\t", data[i][j]))
				}
			}
			if i == 2 {
				imgui.PopStyleColor()
			}
		}
		imgui.PopStyleColor()

	}
	imgui.EndTable()
	imgui.End()
}

func loop() {
	imgui.CreateContext().SetConfigFlagsCurrFrame(imgui.ConfigFlagsNavEnableKeyboard | imgui.ConfigFlagsDockingEnable)
	// io := imgui.CurrentIO()
	// fmt.Printf("A: %v\n", io.Size)
	// showWidgetsDemo()
	imgui.ShowDemoWindow()
	rawMemory()
	disassembledIntcode()
}

func beforeDestroyContext() {
	imgui.PlotDestroyContext()
}

func init() {
	runtime.LockOSThread()
}

func Run(vmI *vm.IntcodeVM) {

	for i := 0; i < 10; i++ {
		barValues = append(barValues, int64(i+1))
	}

	backend, _ = imgui.CreateBackend(imgui.NewGLFWBackend())
	backend.SetBeforeDestroyContextHook(beforeDestroyContext)

	backend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))
	backend.CreateWindow("Hello from cimgui-go", 1200, 900)

	backend.SetDropCallback(func(p []string) {
		fmt.Printf("drop triggered: %v", p)
	})

	backend.SetCloseCallback(func(b imgui.Backend[imgui.GLFWWindowFlags]) {
		fmt.Println("window is closing")
	})

	io := imgui.CurrentIO()
	io.Fonts().AddFontFromFileTTF("./debugger/hack.ttf", 18)

	vmIns = vmI
	backend.Run(loop)

}
