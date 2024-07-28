package debugger

import (
	"fmt"
	"runtime"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/xSaCh/intcode/vm"
)

var (
	vmIns          *vm.IntcodeVM
	showDemoWindow bool
	selected       bool
	backend        imgui.Backend[imgui.GLFWWindowFlags]
	barValues      []int64
	f              bool
)

func showWidgetsDemo() {
	showDemoWindow = true

	did := imgui.DockSpaceOverViewport()

	if showDemoWindow {
		imgui.Begin("Q")
		imgui.Text("HEEELOO")
		imgui.Button("QQQQ")
		imgui.End()
	}

	imgui.SetNextWindowSizeV(imgui.NewVec2(300, 300), imgui.CondOnce)

	imgui.Begin("Window 1")
	if imgui.ButtonV("Click Me", imgui.NewVec2(80, 20)) {
		w, h := backend.DisplaySize()
		fmt.Println(w, h)
	}
	imgui.TextUnformatted("Unformatted text")
	imgui.Checkbox("Show demo window", &showDemoWindow)
	if imgui.BeginCombo("Combo", "Combo preview") {
		imgui.SelectableBoolPtr("Item 1", &selected)
		imgui.SelectableBool("Item 2")
		imgui.SelectableBool("Item 3")
		imgui.EndCombo()
	}

	if imgui.RadioButtonBool("Radio button1", selected) {
		selected = true
	}

	imgui.SameLine()

	if imgui.RadioButtonBool("Radio button2", !selected) {
		selected = false
	}

	// a := imgui.ImNodesGetIO()
	// imgui.Dock

	if f {
		f = false
		// imgui.InternalDockBuilderSetNodeSize(did, imgui.MainViewport().Size())

		mid := did
		// mid = imgui.InternalDockBuilderSplitNode(did, imgui.DirLeft, 0.5, nil, &did)
		rid := imgui.InternalDockBuilderSplitNode(mid, imgui.DirRight, 0.25, nil, &mid)
		imgui.InternalDockBuilderDockWindow("Window 1", mid)
		imgui.InternalDockBuilderDockWindow("Q", rid)
		imgui.InternalDockBuilderFinish(did)
	}

	imgui.End()
}

func rawMemory() {
	data := [][]int{{11101, 3340, 3300, 3300}, {3491, 3949, 2940, 4583}}
	imgui.Begin("table")
	if imgui.BeginTable("table_t", 4) {

		imgui.PushStyleColorU32(imgui.ColText, 0xFF5F5F5F)
		for i := 0; i < 5; i++ {
			imgui.TableNextRow()

			if i == 2 {
				imgui.PushStyleColorU32(imgui.ColText, 0xFFFFFFFF)
			}
			for j := 0; j < 4; j++ {
				imgui.TableSetColumnIndex(int32(j))
				imgui.Text(fmt.Sprintf("%d\t", data[i%2][j]))
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

func parsedIntcode() {
	// data := [][]string{{"ADD #1  #2  #4", "ADD 3 5 #4"}, {"ADD1 #1  #2  #4", "ADD 3 5 #4"},
	// 	{"ADD2 #1  #2  #4", "ADD 3 5 #4"}, {"ADD3 #1  #2  #4", "ADD 3 5 #4"}}
	data, err := GetFormattedMemory(vmIns.Memory)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	imgui.Begin("parsed")
	if imgui.BeginTable("table_t", 4) {

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
	parsedIntcode()
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
	f = true

	io := imgui.CurrentIO()
	fnt := io.Fonts().AddFontFromFileTTF("./debugger/hack.ttf", 18)
	_ = fnt

	vmIns = vmI
	backend.Run(loop)

}
