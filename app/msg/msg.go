package msg

const (
	SUCCESS                          = "ok"
	FAILED                           = "fail"
	ERROR_SIGN_INVALID               = "认证失败"
	ERROR_NEED_LOGIN                 = "用户未登录"
	ERROR_WECHAT_FAILUER             = "微信服务器通讯异常"
	INVALID_PARAMS                   = "请求参数错误"
	ERROR_UPLOAD_SAVE_IMAGE_FAIL     = "保存图片失败"
	ERROR_UPLOAD_CHECK_IMAGE_FAIL    = "检查图片失败"
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT  = "校验图片错误，图片格式或大小有问题"
	ERROR_MOBILE_INVALID             = "手机号格式不正确"
	ERROR_MOBILECODE_NOMATCH         = "验证码不正确"
	ERROR_MOBILE_NOEXIST             = "手机号和验证码不匹配"
	ERROR_MOBILE_SEND_SMS_FAILED     = "短信发送失败,请稍后重试"
	ERROR_TOO_MANY_REQUEST           = "请求太频繁了"
	ERROR_SEARCH_NO_RESULT           = "没有找到相关内容"
	ERROR_COLLECT_BOOK_ADD_OK        = "收藏成功"
	ERROR_COLLECT_BOOK_ADD_FAILED    = "收藏失败"
	ERROR_COLLECT_BOOK_DELETE_FAILED = "取消收藏失败"
	ERROR_COLLECT_BOOK_FRESH_FAILED  = "收藏夹同步失败,请稍后重试"
)
