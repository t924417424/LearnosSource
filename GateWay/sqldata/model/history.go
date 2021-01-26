package model

type History struct {
	Model
	Cid        string `grom:"size:255"`
	BindUser   uint   `json:"-"`
	ImagesName string `grom:"size:255"`
	Status     int    //容器状态 0：创建中，1：Pull镜像文件,2：创建成功，3：创建失败，4：已删除
}
