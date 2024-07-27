package main

import (
	"fmt"
	"runtime"

	imgui "github.com/AllenDang/cimgui-go"
)

var (
	showDemoWindow bool
	value1         int32
	value2         int32
	value3         int32
	values         [2]int32 = [2]int32{value1, value2}
	content        string   = "Let me try"
	r              float32
	g              float32
	b              float32
	a              float32
	color4         [4]float32 = [4]float32{r, g, b, a}
	selected       bool
	backend        imgui.Backend[imgui.GLFWWindowFlags]
	barValues      []int64
	f              bool
)

func callback(data imgui.InputTextCallbackData) int {
	fmt.Println("got call back")
	return 0
}

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

	imgui.InputTextWithHint("Name", "write your name here", &content, 0, callback)
	imgui.Text(content)
	imgui.SliderInt("Slider int", &value3, 0, 100)
	imgui.DragInt("Drag int", &value1)
	imgui.DragInt2("Drag int2", &values)
	value1 = values[0]
	imgui.ColorEdit4("Color Edit3", &color4)

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
	imgui.Begin("parsed")
	data := [][]string{{"ADD #1  #2  #4", "ADD 3 5 #4"}, {"ADD1 #1  #2  #4", "ADD 3 5 #4"},
		{"ADD2 #1  #2  #4", "ADD 3 5 #4"}, {"ADD3 #1  #2  #4", "ADD 3 5 #4"}}

	if imgui.BeginTable("table_t", 2) {

		imgui.PushStyleColorU32(imgui.ColText, 0xFF5F5F5F)
		for i := 0; i < 4; i++ {
			imgui.TableNextRow()

			if i == 2 {
				imgui.PushStyleColorU32(imgui.ColText, 0xFFFFFFFF)
			}
			for j := 0; j < 2; j++ {
				imgui.TableSetColumnIndex(int32(j))
				// if i == 2 {
				// 	imgui.TableSetBgColor(imgui.TableBgTargetCellBg, 0xFFFF4455)
				// }
				imgui.Text(fmt.Sprintf("%s\t", data[i][j]))
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

func main() {

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
	fnt := io.Fonts().AddFontFromFileTTF("hack.ttf", 18)
	_ = fnt
	backend.Run(loop)

}
