package main

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

var colorFg = "#313233"
var colorNorm = "#969696"
var colorWarn = "#E3CA3A"
var colorShit = "#F03434"

// Interval - I just want const.
const Interval = 250

// CPULvl - cpu levels (2 core per symbol)
var CPULvl = [...][5]string{
	{"⠀", "⢀", "⢠", "⢰", "⢸"},
	{"⡀", "⣀", "⣠", "⣰", "⣸"},
	{"⡄", "⣄", "⣤", "⣴", "⣾"},
	{"⡆", "⣆", "⣦", "⣶", "⣾"},
	{"⡇", "⣇", "⣧", "⣷", "⣿"}}

func main() {
	for {
		fmt.Println(`<span font_family="Iosevka" foreground="#969696">` +
			`<span face="Font Awesome 5 Free" weight="normal">&#xf2db; </span>` +
			p() +
			`</span>`)
		time.Sleep(Interval * time.Millisecond)
	}
}

func p() string {
	c, err := cpu.Percent(0, true)
	if err != nil {
		return "⠀⠀⠀⠀"
	}

	cpuOut := ""
	oddCore := 0

	sum := 0.0
	for _, v := range c {
		sum += float64(v)
	}

	common := sum / float64(len(c))
	color := colorNorm

	for i, cpu := range c {
		if cpu > 100 {
			cpu = 100
		}
		grad := int(math.Ceil(cpu / (100 / 4)))

		if cpu > 90 {
			color = colorWarn
		}

		if i%2 == 1 {
			cpuOut += CPULvl[oddCore][grad]
		} else {
			oddCore = grad
		}
	}

	if common >= 50 && common < 90 {
		color = colorWarn
	} else if common >= 90 {
		color = colorShit
	}

	return `<span foreground="` + color + `">` + cpuOut + `</span>`
}
