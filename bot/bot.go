package bot

import (
	"fmt"

	"github.com/riba2534/openwechat"
	"github.com/skip2/go-qrcode"
)

var Bot *openwechat.Bot

func Init() error {
	Bot = openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式
	Bot.UUIDCallback = consoleQrCode                // 注册登陆二维码回调
	reloadStorage := openwechat.NewFileHotReloadStorage("token.json")
	return Bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())
}

func consoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.High)
	fmt.Println(q.ToSmallString(false))
}
