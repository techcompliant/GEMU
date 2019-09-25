package gemu

type MultiStorage struct {
	storage []Storage
}

func NewMultiStorage(instorage ...Storage) Storage {
	return &MultiStorage{storage: instorage}
}

func (MS *MultiStorage) Exists(Item string) bool {
	for _, S := range MS.storage {
		if S.Exists(Item) {
			return true
		}
	}
	return false
}

func (MS *MultiStorage) Length(Item string) int {
	for _, S := range MS.storage {
		if S.Exists(Item) {
			return S.Length(Item)
		}
	}
	return 0
}

func (MS *MultiStorage) Read(Item string, offset int, data []byte) {
	for _, S := range MS.storage {
		if S.Exists(Item) {
			S.Read(Item, offset, data)
			return
		}
	}
}

func (MS *MultiStorage) Write(Item string, offset int, data []byte) {
	for _, S := range MS.storage {
		if S.Exists(Item) {
			S.Write(Item, offset, data)
			return
		}
	}
	MS.storage[0].Write(Item, offset, data)
}
