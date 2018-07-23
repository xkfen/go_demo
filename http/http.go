package http

import "gopkg.in/gin-gonic/gin.v1"

type BaseResp struct {
	// 成功标识
	Success bool	`json:"success"`
	// 信息
	Info string	`json:"info"`
	// 错误信息
	ErrorMsg string `json:"errmsg"`
}

// 返回成功
func RenderSuccess(data interface{}, c *gin.Context) {
	c.JSON(200, data)
}

// 返回失败（标准化返回）
func Render400(info string, c *gin.Context) {
	c.JSON(400, BaseResp{Success: false, Info: info})
}

// 返回权限验证失败
func Render403(info string, c *gin.Context) {
	c.JSON(403, BaseResp{Success: false, Info: info})
}