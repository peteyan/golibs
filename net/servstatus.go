package net

// ----------------------------------------------------------------
// 参照http请求结构体：src/net/http/status
// StatusCode定义：1位状态code字母+4位code数字
// 状态code：P(PASS), E(ERROR), U(UNKNOWN)
// ----------------------------------------------------------------

const (
	// StatusCodeOK OK
	StatusCodeOK = "P0200"

	// StatusCodeErrorParams 参数错误
	StatusCodeErrorParams = "E0101"

	// StatusCodeErrorDataNotFound 数据未找到
	StatusCodeErrorDataNotFound = "E0102"

	// StatusCodeErrorInternal 服务内部错误
	StatusCodeErrorInternal = "E0103"
)

// StatusMsg 根据状态码和自定义信息返回整个消息，注意：第一个参数必须要给状态码
//
// eg：StatusMsg(StatusCodeErrorParams, "-xx参数值不能为空")
func StatusMsg(msg ...string) string {
	s := ""
	if len(msg) == 0 {
		return s
	}
	switch msg[0] {
	case StatusCodeOK:
		s += "OK"
		break
	case StatusCodeErrorParams:
		s += "参数错误"
		break
	case StatusCodeErrorDataNotFound:
		s += "数据未找到"
		break
	case StatusCodeErrorInternal:
		s += "服务内部错误"
		break
	default:
		break
	}
	for i, v := range msg {
		if i == 0 {
			continue
		}
		s += v
	}
	return s
}
