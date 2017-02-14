package gemu // import "gemu.techcompliant.com/gemu"

type FlipStorage struct {
	Storage
}

func NewFlipStorage(base Storage) Storage {
	return &FlipStorage{base}
}

func (FS *FlipStorage) Read(Item string, offset int, data []byte) {
	FS.Storage.Read(Item, offset, data)
	for l1 := 0; l1 < len(data)/2; l1++ {
		data[l1*2], data[l1*2+1] = data[l1*2+1], data[l1*2]
	}
}

func (FS *FlipStorage) Write(Item string, offset int, data []byte) {
	for l1 := 0; l1 < len(data)/2; l1++ {
		data[l1*2], data[l1*2+1] = data[l1*2+1], data[l1*2]
	}
	FS.Write(Item, offset, data)
}
