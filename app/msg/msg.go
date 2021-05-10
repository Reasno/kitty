package msg

const (
	ErrorNeedLogin                 = "用户未登录"
	ErrorLogin                     = "登陆失败"
	ErrorCorruptedData             = "数据被污染"
	ErrorWechatFailure             = "微信服务器通讯异常"
	ErrorWechatLogin               = "微信登陆失败"
	ErrorUpload                    = "无法上传图片"
	ErrorAlreadyBind               = "已经被其他账号绑定"
	ErrorTooFrequent               = "请求太频繁了"
	ErrorMissingOpenid             = "腾讯服务器验证未通过"
	ErrorGetCode                   = "获取验证码异常"
	ErrorSendCode                  = "生成验证码异常"
	ErrorDatabaseFailure           = "数据库异常"
	InvalidParams                  = "请求参数错误"
	ErrorMobileCode                = "手机号和验证码不匹配"
	ErrorJwtFailure                = "无法生成签名"
	ErrorRecordNotFound            = "目标不存在"
	ErrorCircledInvitation         = "对方不能被邀请，试试邀请其他人吧"
	ErrorRelationAlreadyExists     = "已经邀请过了"
	ErrorExtraNotFound             = "所请求的信息不存在或已过期"
	WxSuccess                      = "微信用户%d成功登录"
	MobileSuccess                  = "手机用户%d成功登录"
	DeviceSuccess                  = "设备用户%d成功登录"
	RewardClaimed                  = "奖励已经被领取过了"
	OrientationHasNotBeenCompleted = "初始任务还没有完成"
	RepeatedInviteCode             = "邀请码只能提交一次，不可修改"
	InvalidInviteCode              = "不合法的邀请码"
	NoRewardAvailable              = "没有可供领取的奖励"
	InvalidInviteTarget            = "不合法的邀请对象"
	InvalidInviteSequence          = "邀请者的注册日期晚于被邀请者"
	XTastAbnormally                = "任务服务繁忙，请稍后再试"
	ReenteringCode                 = "不能重复填写邀请码"
	AdminOnly                      = "当前操作需要管理员权限"
	AlreadyDeleted                 = "用户已被删除"
	ServerBug                      = "应用程序异常"
)
