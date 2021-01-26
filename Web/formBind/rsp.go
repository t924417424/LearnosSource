package formBind

type StatusCode int

const (
	FormErr     StatusCode = 101
	Success     StatusCode = 200
	TokenErr    StatusCode = 301
	TokenExpErr StatusCode = 302
	GateWayErr  StatusCode = 304
)

type Rsp struct {
	Code  StatusCode  `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Token string      `json:"token"`
}
