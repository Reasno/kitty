package msg

const (
	ErrorNeedLogin       = "用户未登录"
	ErrorLogin           = "登陆失败"
	ErrorWechatFailure   = "微信服务器通讯异常"
	ErrorWechatLogin     = "微信登陆失败"
	ErrorUpload          = "无法上传图片"
	ErrorAlreadyBind     = "已经被其他账号绑定"
	ErrorTooFrequent     = "请求太频繁了"
	ErrorMissingOpenid   = "腾讯服务器验证未通过"
	ErrorGetCode         = "获取验证码异常"
	ErrorSendCode        = "生成验证码异常"
	ErrorDatabaseFailure = "数据库异常"
	InvalidParams        = "请求参数错误"
	ErrorMobileCode      = "手机号和验证码不匹配"
	ErrorJwtFailure      = "无法生成签名"
	ErrorRecordNotFound  = "目标不存在"
	ErrorExtraNotFound   = "所请求的信息不存在或已过期"
	WxSuccess            = "微信用户%d成功登录"
	MobileSuccess        = "手机用户%d成功登录"
	DeviceSuccess        = "设备用户%d成功登录"
)
