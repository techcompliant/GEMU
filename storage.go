package gemu

var defaultStorage Storage

func SetStorage(storage Storage) {
	defaultStorage = storage
}

type Storage interface {
	Exists(Item string) bool
	Length(Item string) int
	Read(Item string, offset int, data []byte)
	Write(Item string, offset int, data []byte)
}
