package msg

const (
	ErrorNeedLogin       = "用户未登录"
	ErrorLogin = "登陆失败"
	ErrorWechatFailure   = "微信服务器通讯异常"
	ErrorUpload          = "无法上传图片"
	ErrorMissingOpenid   = "OpenID缺失"
	ErrorDatabaseFailure = "数据库异常"
	InvalidParams        = "请求参数错误"
	ErrorMobileCode      = "手机号和验证码不匹配"
	ErrorJwtFailure      = "无法生成签名"
)
