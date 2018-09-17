package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/net"
)

var minC uint64 = 120
var maxC uint64 = 254
var bColor = "#969696"
var kbColor = "#aaaaaa"
var mbColor = "#cccccc"

// Interval - I just want const.
const Interval = 1000

// Interface name.
const Interface = "wlp13s0b1"

// MaxSpeed hm...
const MaxSpeed = 2000000

func main() {
	for {
		fmt.Println(`<span size="7600">` + n() + `</span>`)
		time.Sleep(Interval * time.Millisecond)
	}
}

var bs uint64
var br uint64
var ts = 1.0
var tr = 1.0
var ss float64
var sr float64
var ssPrev float64
var srPrev float64
var bytesSent uint64
var bytesRecv uint64
var bsOut string
var brOut string

func n() string {
	n, err := net.IOCounters(true)
	if err != nil {
		return ""
	}

	var stat *net.IOCountersStat
	for _, c := range n {
		if c.Name == Interface {
			stat = &c
			break
		}
	}
	if stat == nil {
		return ""
	}

	if stat.BytesSent != bytesSent {
		bs = stat.BytesSent - bytesSent
		ss = float64(bs) / ts
		ss = ssPrev*0.1 + ss*0.9
		ts = float64(Interval) / 1000
	} else {
		ss /= 1.1
		ts += (float64(Interval) / 1000)
	}
	if stat.BytesRecv != bytesRecv {
		br = stat.BytesRecv - bytesRecv
		sr = float64(br) / tr
		sr = srPrev*0.1 + sr*0.9
		tr = float64(Interval) / 1000
	} else {
		sr /= 1.1
		tr += (float64(Interval) / 1000)
	}

	lss := ss
	if ss > MaxSpeed {
		lss = MaxSpeed
	}
	lsr := sr
	if sr > MaxSpeed {
		lsr = MaxSpeed
	}
	csBase := uint64(lss/MaxSpeed*float64(maxC-minC)) + minC
	crBase := uint64(lsr/MaxSpeed*float64(maxC-minC)) + minC
	cs := "#" + strings.Repeat(strconv.FormatUint(csBase, 16), 3)
	cr := "#" + strings.Repeat(strconv.FormatUint(crBase, 16), 3)

	bytesSent = stat.BytesSent
	bytesRecv = stat.BytesRecv
	ssPrev = ss
	srPrev = sr

	shortS, us := shortSize(ss)
	shortR, ur := shortSize(sr)

	return `<span face="Meslo LG S" size="7600" foreground="` + cs + `">` +
		fmt.Sprintf("%3v", shortS) + us +
		`</span>` +

		`<span face="Font Awesome 5 Free" weight="normal" foreground="` + cs + `"> &#xf062;</span>` +

		`<span face="Meslo LG S" size="7600" foreground="` + cr + `">` +
		fmt.Sprintf("%4v", shortR) + ur +
		`</span>` +

		`<span face="Font Awesome 5 Free" weight="normal" foreground="` + cr + `"> &#xf063;</span>`
}

func shortSize(bytes float64) (string, string) {
	if bytes <= 1024 {
		if bytes >= 10 {
			return strconv.FormatFloat(bytes, 'f', 0, 64), "b"
		}
		return strconv.FormatFloat(bytes, 'f', 1, 64), "b"
	}

	if kb := bytes / 1024; kb <= 1024 {
		if kb >= 10 {
			return strconv.FormatFloat(kb, 'f', 0, 64), "k"
		}
		return strconv.FormatFloat(kb, 'f', 1, 64), "k"
	}

	if mb := bytes / 1048576; mb <= 1024 {
		if mb >= 10 {
			return strconv.FormatFloat(mb, 'f', 0, 64), "m"
		}
		return strconv.FormatFloat(mb, 'f', 1, 64), "m"
	}

	gb := bytes / 1073741824
	if gb >= 10 {
		return strconv.FormatFloat(gb, 'f', 0, 64), "g"
	}
	return strconv.FormatFloat(gb, 'f', 1, 64), "g"
}
