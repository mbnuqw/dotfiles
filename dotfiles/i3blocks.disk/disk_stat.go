package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/disk"
)

// Disk stats.
type Disk struct {
	Dev   string
	Mnt   string
	Name  string
	Usage *disk.UsageStat
	Stats disk.IOCountersStat
	PrevR uint64
	PrevW uint64
}

var disks = make([]Disk, 0, 5)

var levels = [...]string{"⠀",
	"⡀", "⣀",
	"⣄", "⣤",
	"⣦", "⣶",
	"⣷", "⣿"}

func main() {
	for {
		fmt.Println(`<span font_family="Iosevka" foreground="#ffffff">` +
			d() +
			`</span>`)
		time.Sleep(1000 * time.Millisecond)
	}
}

func d() string {
	parts := GetPartitions()
	disks = UpdateStats(parts)

	var output string
	for i := range disks {
		d := &disks[i]

		if i > 0 {
			output += "  "
		}

		if d.Mnt == "/" {
			output += `<span face="Font Awesome 5 Free" foreground="#cecece" weight="bold">&#xf0a0;</span>`
		} else {
			output += `<span face="Font Awesome 5 Free" foreground="#cecece" weight="normal">&#xf0a0;</span>`
		}

		readb := d.Stats.ReadBytes - d.PrevR
		writeb := d.Stats.WriteBytes - d.PrevW
		rv, ru := shortSize(float64(readb))
		wv, wu := shortSize(float64(writeb))

		output += `<span face="Meslo LG S" size="7600">` +
			fmt.Sprintf("%4v", rv) + ru +
			`</span>`

		output += `<span face="Meslo LG S" size="7600">` +
			fmt.Sprintf("%4v", wv) + wu +
			`</span>`

		d.PrevR = d.Stats.ReadBytes
		d.PrevW = d.Stats.WriteBytes
	}

	return output
}

// GetPartitions parttttttttt
func GetPartitions() []disk.PartitionStat {
	parts, err := disk.Partitions(false)
	if err != nil {
		log.Fatal("Oooufffu...", err)
	}
	return parts
}

// UpdateStats sss
func UpdateStats(parts []disk.PartitionStat) []Disk {
	var devs = make([]string, 0, 5)
	disksLen := len(disks)
	partsLen := len(parts)
	count := int(math.Max(float64(disksLen), float64(partsLen)))

	for i := 0; i < count; i++ {
		// Update
		if partsLen > i && disksLen > i {
			part := &parts[i]
			usage, err := disk.Usage(part.Mountpoint)
			if err != nil {
				log.Fatal("Ophu", err)
			}

			disks[i].Dev = part.Device
			disks[i].Mnt = part.Mountpoint
			disks[i].Usage = usage

			devs = append(devs, part.Device)
		}

		// Add new
		if partsLen > i && disksLen <= i {
			part := &parts[i]
			usage, err := disk.Usage(part.Mountpoint)
			if err != nil {
				log.Fatal("Ophu", err)
			}

			disks = append(disks, Disk{
				Dev:   part.Device,
				Mnt:   part.Mountpoint,
				Usage: usage,
			})

			devs = append(devs, part.Device)
		}

		// Remove
		if partsLen <= i && disksLen > i {
			disks = append(disks[:i], disks[i+1:]...)
		}
	}

	counts, err := disk.IOCounters(devs...)
	if err != nil {
		log.Fatal("Offu", err)
	}

	for name, stat := range counts {
		for i := range disks {
			disk := &disks[i]
			if strings.Contains(disk.Dev, name) {
				disk.Name = name
				disk.Stats = stat
			}
		}
	}

	return disks
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
