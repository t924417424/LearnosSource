package formBind

type Register struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Phone string `form:"phone" binding:"required"`
	Verify string `form:"verify" binding:"required"`
}

type SendCode struct {
	Phone string `form:"phone" binding:"required"`
}

type Login struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type ContainerApi struct {
	Cid string `form:"cid" binding:"required"`
	Limit bool `form:"limit"`
}

type CreateContainer struct {
	Id uint32 `form:"id" binding:"required"`
}