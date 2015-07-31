package hiapns

import "errors"

var (
	ErrBadValue     = errors.New("值不合法")
	ErrNoClient     = errors.New("client 不存在")
	ErrCertMissing  = errors.New("cert 不能为空")
	ErrKeyMissing   = errors.New("key 不能为空")
	ErrAlertMissing = errors.New("alert 不能为空")
	ErrBadAlert     = errors.New("alert 只能是字符串")
	ErrBadBadge     = errors.New("badge 只能是数字")
	ErrBadSound     = errors.New("sound 只能是字符串")
	ErrBadFormat    = errors.New("参数格式不正确")
)
