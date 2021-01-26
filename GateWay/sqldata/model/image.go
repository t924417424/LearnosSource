package model

type Image struct {
	Model
	Logo           string `grom:"size:255"`
	Name           string `grom:"size:255"`
	Cmd            string `grom:"size:255"`
	Network        bool   `gorm:"default:false"`
	Memory         int64  `gorm:"type:bigint;default:0"`
	NetWorkIoLimit uint64 `gorm:"type:bigint;default:0"`
	BlockIoLimit   uint64 `gorm:"type:bigint;default:0"`
}
