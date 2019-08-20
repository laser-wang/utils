package captcha

import (
	"caton/xw.wang/utils"
	"caton/xw.wang/utils/cache"
	"github.com/dchest/captcha"
)

var (
	SR *StoreRedis
)

type StoreRedis struct {
}
type ImgBytes struct {
	Img []byte
}

func (s *StoreRedis) Set(id string, digits []byte) {
	obj := new(ImgBytes)
	obj.Img = digits
	cache.DoSet(id, obj, utils.TIME_MINUTE_FIVE)
}
func (s *StoreRedis) Get(id string, clear bool) (digits []byte) {
	obj := new(ImgBytes)
	_ = cache.DoGet(id, obj)
	return obj.Img
}

func InitCaptcha() {
	SR = new(StoreRedis)
	captcha.SetCustomStore(SR)
}

func NewLen(length int) (id string) {
	captchaId := captcha.NewLen(length)
	return captchaId
}

func VerifyString(id string, digits string) bool {
	ret := captcha.VerifyString(id, digits)
	// 验证一次以后就失效
	cache.DoDel(id)
	return ret
}

// 验证码不失效，允许再次发送
func VerifyStringNODel(id string, digits string) bool {
	ret := captcha.VerifyString(id, digits)
	return ret
}
