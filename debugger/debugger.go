package debugger

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/xSaCh/intcode/vm"
)

var (
	vmIns      *vm.IntcodeVM
	backend    imgui.Backend[imgui.GLFWWindowFlags]
	barValues  []int64
	outputLog  []int
	canStop    bool
	canSetDock bool
	isAscii    bool
	addNewLine bool   = true
	speed      int32  = 200
	lblRun     string = "Run"
)

func rawMemory() {
	memMaxWidth := len(strconv.Itoa(len(vmIns.Memory)))
	imgui.Begin("Raw Memory")
	if imgui.BeginTable("tbl_raw_mem", 2) {

		imgui.PushStyleColorU32(imgui.ColText, 0xFF5F5F5F)
		for i := 0; i < len(vmIns.Memory); i++ {
			imgui.TableNextRow()

			if i == vmIns.PcRegister {
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
			if i == vmIns.PcRegister {
				imgui.PopStyleColor()
			}
		}
		imgui.PopStyleColor()

		imgui.EndTable()
	}
	imgui.End()
}

func disassembledIntcode() {
	data, ind, err := GetFormattedMemory(vmIns.Memory, vmIns.PcRegister)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	// fmt.Printf("ind: %v\n", ind)
	imgui.Begin("Disassembled Intcode")
	if imgui.BeginTable("tbl_dis_int", 4) {

		imgui.PushStyleColorU32(imgui.ColText, 0xFF5F5F5F)
		for i := 0; i < len(data); i++ {
			imgui.TableNextRow()

			if i == ind {
				imgui.PushStyleColorU32(imgui.ColText, 0xFFFFFFFF)
			}
			for j := 0; j < 4; j++ {
				imgui.TableSetColumnIndex(int32(j))

				if j < len(data[i]) {
					imgui.Text(fmt.Sprintf("%s\t", data[i][j]))
				}
			}
			if i == ind {
				imgui.PopStyleColor()
			}
		}
		imgui.PopStyleColor()

		imgui.EndTable()
	}
	imgui.End()
}

func console() {

	imgui.Begin("Console")
	imgui.Checkbox("ASCII", &isAscii)
	imgui.SameLine()
	imgui.Checkbox("Add NewLine", &addNewLine)

	output := ""
	if !addNewLine {
		for _, v := range outputLog {
			if isAscii {
				output += fmt.Sprintf("%c", v)
			} else {
				output += fmt.Sprintf("%d", v)
			}
		}
		imgui.Text(output)
	} else {
		for _, v := range outputLog {
			if isAscii {
				imgui.Text(fmt.Sprintf("%c", v))
			} else {
				imgui.Text(fmt.Sprintf("%d", v))
			}
		}
	}
	imgui.End()
}

func loop() {

	// imgui.ShowDemoWindow()

	//Raw Memory
	//Disassembled Intcode
	imgui.CreateContext().SetConfigFlagsCurrFrame(imgui.ConfigFlagsNavEnableKeyboard | imgui.ConfigFlagsDockingEnable)
	did := imgui.DockSpaceOverViewport()

	rawMemory()
	disassembledIntcode()

	imgui.Begin("Controller")
	if imgui.Button("Step") {
		if !canStop {
			canStop, _ = vmIns.Step()
		}
	}

	if imgui.Button("Reset") {
		lblRun = "Run"
		vmIns.LoadProgram(slices.Clone(vmIns.InitMemory))
		outputLog = []int{}
		canStop = false

	}

	if imgui.Button(lblRun) {
		if lblRun == "Pause" {
			lblRun = "Run"
		} else {
			lblRun = "Pause"
		}
	}

	imgui.SliderInt("Execution Speed (in ms)", &speed, 0, 1000)
	imgui.End()

	console()
	if canSetDock {
		canSetDock = false

		imgui.InternalDockBuilderSetNodeSize(did, imgui.MainViewport().Size())
		_ = did
		mid := imgui.InternalDockBuilderSplitNode(did, imgui.DirLeft, 0.7, nil, &did)
		rid := imgui.InternalDockBuilderSplitNode(mid, imgui.DirLeft, 0.3, nil, &mid)
		imgui.InternalDockBuilderDockWindow("Raw Memory", rid)
		imgui.InternalDockBuilderDockWindow("Disassembled Intcode", mid)

		cid := did
		cbid := imgui.InternalDockBuilderSplitNode(cid, imgui.DirDown, 0.3, nil, &cid)
		imgui.InternalDockBuilderDockWindow("Controller", cid)
		imgui.InternalDockBuilderDockWindow("Console", cbid)
		imgui.InternalDockBuilderFinish(did)
	}

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
	backend.CreateWindow("Intcode Visualizer", 1200, 900)

	backend.SetDropCallback(func(p []string) {
		fmt.Printf("drop triggered: %v", p)
	})

	backend.SetCloseCallback(func(b imgui.Backend[imgui.GLFWWindowFlags]) {
		fmt.Println("window is closing")
	})

	io := imgui.CurrentIO()
	io.Fonts().AddFontFromFileTTF("./debugger/hack.ttf", 18)

	vmIns = vmI

	vmIns.OutputFunc = func(i int) {
		outputLog = append(outputLog, i)
	}

	canSetDock = true

	go func() {
		for {
			if lblRun == "Pause" && !canStop {
				canStop, _ = vmIns.Step()
				backend.Refresh()
			}
			time.Sleep(time.Millisecond * time.Duration(speed))
		}
	}()

	backend.Run(loop)
}
