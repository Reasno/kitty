package msg

const (
	ErrorNeedLogin       = "用户未登录"
	ErrorLogin           = "登陆失败"
	ErrorWechatFailure   = "微信服务器通讯异常"
	ErrorWechatLogin     = "微信登陆失败"
	ErrorUpload          = "无法上传图片"
	ErrorTooFrequent     = "请求太频繁了"
	ErrorMissingOpenid   = "OpenID缺失"
	ErrorGetCode         = "获取验证码异常"
	ErrorSendCode        = "生成验证码异常"
	ErrorDatabaseFailure = "数据库异常"
	InvalidParams        = "请求参数错误"
	ErrorMobileCode      = "手机号和验证码不匹配"
	ErrorJwtFailure      = "无法生成签名"
	ErrorUserNotFound    = "目标用户不存在"
	WxSuccess            = "微信用户%d成功登录"
	MobileSuccess        = "手机用户%d成功登录"
	DeviceSuccess        = "设备用户%d成功登录"
)
