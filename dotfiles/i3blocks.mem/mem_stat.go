package main

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/mem"
)

var colorFg = "#313233"
var colorNorm = "#ffffff"
var colorWarn = "#E3CA3A"
var colorShit = "#F03434"

// Interval - I just want const.
const Interval = 1500

// MemLvl - memory levels.
var MemLvl = [...]string{"⠀⠀",
	"⡀⠀", "⣀⠀", "⣀⡀", "⣀⣀",
	"⣄⣀", "⣤⣀", "⣤⣄", "⣤⣤",
	"⣦⣤", "⣶⣤", "⣶⣦", "⣶⣶",
	"⣷⣶", "⣿⣶", "⣿⣷", "⣿⣿"}

func main() {
	for {
		fmt.Println(`<span font_family="Iosevka">` +
			// `<span face="Font Awesome 5 Free" weight="normal">&#xf1c0; </span>` +
			m() +
			`</span>`)
		time.Sleep(Interval * time.Millisecond)
	}
}

func m() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "⠀⠀"
	}

	swapInfo, err := mem.SwapMemory()
	if err != nil {
		return "  "
	}

	memUsed := memInfo.UsedPercent
	if memUsed > 100 {
		memUsed = 100
	}
	memGrad := int(math.Ceil(memUsed / (100 / (float64(len(MemLvl)) - 1))))

	color := colorNorm
	if swapInfo.Used > 0 && swapInfo.Used < 0xffffff {
		color = colorWarn
	} else if swapInfo.Used >= 0xffffff {
		color = colorShit
	}

	return `<span foreground="` + color + `">` + MemLvl[memGrad] + `</span>`
}
