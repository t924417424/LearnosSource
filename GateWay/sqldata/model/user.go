package model

type User struct {
	Model
	Username     string    `grom:"not null;size:50"`
	Password     string    `grom:"not null;size:32"`
	Phone       string    `gorm:"not null;unique;size:11"`
	UseList      int       `gorm:"not null;default:1"`
	MyContainers []History `gorm:"ForeignKey:BindUser"`
}
