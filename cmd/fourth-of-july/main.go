package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/Andrew-Wichmann/asciiphysics"
	"github.com/Andrew-Wichmann/fourth-of-july/pkg/firework"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const INIT state = 0
const RUNNING state = 1
const PAUSED state = 2

type app struct {
	state         state
	canvas        asciiphysics.Canvas
	demo_firework firework.Model
}

func (app) Init() tea.Cmd {
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
	}
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		a.canvas = asciiphysics.NewCanvas(msg.Width, msg.Height*2)
		a.demo_firework = firework.New(asciiphysics.Vector{X: float64(msg.Width) / 2, Y: float64(msg.Height) / 2})
		a.canvas.AddDrawable(a.demo_firework)
		a.state = RUNNING
		cmds = append(cmds, a.canvas.Init())
	}
	if a.state == RUNNING {
		canvas, cmd := a.canvas.Update(msg)
		a.canvas = canvas
		cmds = append(cmds, cmd)
	}
	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	if a.state == INIT {
		return ""
	} else if a.state == RUNNING {
		return a.canvas.View()
	} else if a.state == PAUSED {
		return "paused"
	}
	panic(fmt.Sprintf("unknown application state: %d", a.state))
}

func main() {
	prog := tea.NewProgram(app{})
	if os.Getenv("ENABLE_PROFILING") == "true" {
		// cpu profile
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()

		// heap profile
		f, err = os.Create("heap.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(err)
		}
		go func() {
			time.Sleep(5 * time.Second)
			prog.Send(tea.Quit())
		}()
	}
	_, err := prog.Run()
	if err != nil {
		panic(err)
	}

}
