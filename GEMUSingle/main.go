package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"gemu.techcompliant.com/gemu"
	"github.com/andyleap/tinyfb"
)

func GetColor(c uint16) color.Color {
	return color.RGBA{
		R: uint8((c >> 8 & 0x0f) << 4),
		G: uint8((c >> 4 & 0x0f) << 4),
		B: uint8((c >> 0 & 0x0f) << 4),
	}
}

var RomImage = flag.String("rom", "internal/bbos.bin", "Filename of rom image to use (internal bbos by default)")
var RomFlip = flag.Bool("noromflip", false, "Don't endian flip the rom")

type FloppyImages []string

func (fi *FloppyImages) String() string {
	return fmt.Sprintf("%s", *fi)
}

func (fi *FloppyImages) Set(value string) error {
	*fi = append(*fi, value)
	return nil
}

func main() {
	fis := &FloppyImages{}
	flag.Var(fis, "floppy", "Floppy images to use (can specify multiple times for multiple drives)")

	flag.Parse()

	t := tinyfb.New("DCPU", (128+12)*4, (96+12)*4)
	go func() {
		t.Run()
		os.Exit(0)
	}()

	gemu.SetStorage(gemu.NewMultiStorage(AssetStorage{Root: "internal/"}, gemu.NewDiskStorage(".")))

	cpu := gemu.NewDCPU(0)

	rom := gemu.NewRom(*RomImage, !*RomFlip)
	cpu.Attach(rom)
	rom.SetUp(cpu)

	clock := gemu.NewClock()
	cpu.Attach(clock)
	clock.SetUp(cpu)

	lem := gemu.NewLem1802()
	cpu.Attach(lem)
	lem.SetUp(cpu)

	keyboard := gemu.NewKeyboard()
	cpu.Attach(keyboard)
	keyboard.SetUp(cpu)

	floppies := []*gemu.M35FD{}

	for _, fi := range *fis {
		floppy := gemu.NewM35FD(true)
		cpu.Attach(floppy)
		floppy.SetUp(cpu)
		floppy.ChangeDisk(fi)
		floppies = append(floppies, floppy)
	}

	cpu.Start()

	lemImageBig := image.NewRGBA(image.Rect(0, 0, (128+12)*4, (96+12)*4))

	t.Char(func(char string, mods int) {
		switch char {
		case "BackSpace", "\x08":
			char = "\x10"
		case "Return", "\x0d":
			char = "\x11"
		case "Insert":
			char = "\x12"
		case "Delete":
			char = "\x13"
		case "Up":
			char = "\x80"
		case "Down":
			char = "\x81"
		case "Left":
			char = "\x82"
		case "Right":
			char = "\x83"
		}

		if len(char) == 1 {
			keyboard.ParsedKey(uint16(char[0]))
		} else {
			//log.Println("Got: ", char)
		}
	})

	t.Key(func(key string, mods int, press bool) {
		switch key {
		case "BackSpace", "\x08":
			key = "\x10"
		case "Return", "\x0d":
			key = "\x11"
		case "Insert":
			key = "\x12"
		case "Delete":
			key = "\x13"
		}
		keyboard.RawKey(uint16(key[0]), press)
	})

	go func() {
		for {
			time.Sleep(20 * time.Millisecond)

			if lem.DspMem == 0 {
				continue
			}
			dspRam := lem.GetMem().GetRaw()[lem.DspMem:]
			fontRam := gemu.LemDefFont
			if lem.FontMem != 0 {
				fontRam = lem.GetMem().GetRaw()[lem.FontMem:]
			}
			palRam := gemu.LemDefPal
			if lem.PalMem != 0 {
				palRam = lem.GetMem().GetRaw()[lem.PalMem:]
			}
			cl := GetColor(palRam[lem.Border&0xf])
			for x := 0; x < (128+12)*4; x++ {
				for y := 0; y < 6*4; y++ {
					lemImageBig.Set(x, y, cl)
					lemImageBig.Set(x, ((96+12)*4)-y, cl)
				}
			}
			for x := 0; x < (6)*4; x++ {
				for y := 0; y < (96+12)*4; y++ {
					lemImageBig.Set(x, y, cl)
					lemImageBig.Set((128+12)*4-x, y, cl)
				}
			}

			for ty := 0; ty < 12; ty++ {
				vtb := uint16(0x0100)
				vta := uint16(0x0001)
				for y := 0; y < 8; y++ {
					for x := 0; x < 32; x++ {
						vtw := dspRam[x+ty*32]
						clb := palRam[(vtw>>8)&0x0f]
						clf := palRam[(vtw>>12)&0x0f]
						vtw &= 0x7f
						vtw <<= 1
						f := fontRam[vtw]
						cl := clf
						if f&vtb == 0 {
							cl = clb
						}
						for sx := 0; sx < 4; sx++ {
							for sy := 0; sy < 4; sy++ {
								lemImageBig.Set(((x*4)+0)*4+sx+6*4, (ty*8+y)*4+sy+6*4, GetColor(cl))
							}
						}

						f = fontRam[vtw]
						cl = clf
						if f&vta == 0 {
							cl = clb
						}
						for sx := 0; sx < 4; sx++ {
							for sy := 0; sy < 4; sy++ {
								lemImageBig.Set(((x*4)+1)*4+sx+6*4, (ty*8+y)*4+sy+6*4, GetColor(cl))
							}
						}

						f = fontRam[vtw+1]
						cl = clf
						if f&vtb == 0 {
							cl = clb
						}
						for sx := 0; sx < 4; sx++ {
							for sy := 0; sy < 4; sy++ {
								lemImageBig.Set(((x*4)+2)*4+sx+6*4, (ty*8+y)*4+sy+6*4, GetColor(cl))
							}
						}

						f = fontRam[vtw+1]
						cl = clf
						if f&vta == 0 {
							cl = clb
						}
						for sx := 0; sx < 4; sx++ {
							for sy := 0; sy < 4; sy++ {
								lemImageBig.Set(((x*4)+3)*4+sx+6*4, (ty*8+y)*4+sy+6*4, GetColor(cl))
							}
						}

					}
					vtb <<= 1
					vta <<= 1
				}
			}
			t.Update(lemImageBig)
		}
	}()

	for {
		cpu.Tick(1)
		clock.Tick(1)
		for _, fd := range floppies {
			fd.Tick(1)
		}
	}
}

type AssetStorage struct {
	Root string
}

func (a AssetStorage) Exists(Item string) bool {
	if strings.HasPrefix(Item, a.Root) {
		Item = strings.TrimPrefix(Item, a.Root)
	} else {
		return false
	}
	_, err := Asset(Item)
	if err != nil {
		return false
	}
	return true
}
func (a AssetStorage) Length(Item string) int {
	if strings.HasPrefix(Item, a.Root) {
		Item = strings.TrimPrefix(Item, a.Root)
	} else {
		return 0
	}
	asset, err := AssetInfo(Item)
	if err != nil {
		return 0
	}
	return int(asset.Size())
}
func (a AssetStorage) Read(Item string, offset int, data []byte) {
	if strings.HasPrefix(Item, a.Root) {
		Item = strings.TrimPrefix(Item, a.Root)
	} else {
		return
	}
	asset, err := Asset(Item)
	if err != nil {
		return
	}
	copy(data, asset[offset:])
}
func (a AssetStorage) Write(Item string, offset int, data []byte) {
	if strings.HasPrefix(Item, a.Root) {

	} else {
		return
	}
	log.Println("Writing to internal assets is not supported!")
}
