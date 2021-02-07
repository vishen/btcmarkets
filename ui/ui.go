package userinterface

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func Run() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

	table3 := widgets.NewTable()
	table3.Rows = [][]string{
		[]string{"Market", "BestBid", "BestAsk", "LastPrice", "5m", "10m"},
		[]string{"AAA", "BBB", "CCC"},
	}
	table3.TextStyle = ui.NewStyle(ui.ColorWhite)
	table3.RowSeparator = true
	table3.BorderStyle = ui.NewStyle(ui.ColorGreen)
	table3.SetRect(0, 0, termWidth, termHeight)
	table3.FillRow = true
	table3.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)

	ui.Render(table3)

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "h":
				table3.Rows = [][]string{
					[]string{"Market", "BestBid", "BestAsk", "LastPrice", "5m", "10m"},
					[]string{"CCC", "DDDD", "EEE"},
				}
				ui.Render(table3)
			case "q", "<C-c>":
				return
			}
		}
	}
}
