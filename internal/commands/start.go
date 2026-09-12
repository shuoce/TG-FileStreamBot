package commands

import (
        "fmt"
	"EverythingSuckz/fsb/internal/authorization"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/tg"
)

func (m *command) LoadStart(dispatcher dispatcher.Dispatcher) {
	log := m.log.Named("start")
	defer log.Sugar().Info("Loaded")
	dispatcher.AddHandler(handlers.NewCommand("start", start))
}

func start(ctx *ext.Context, u *ext.Update) error {
        fmt.Println("========== START COMMAND RECEIVED ==========")
	chatId := u.EffectiveChat().GetID()
	chatId := u.EffectiveChat().GetID()
	peerChatId := ctx.PeerStorage.GetPeerById(chatId)

	if peerChatId.Type != int(storage.TypeUser) {
		return dispatcher.EndGroups
	}

	if !authorization.IsAuthorized(chatId) {
		_, err := ctx.Reply(
			u,
			ext.ReplyTextString("🔐 暂无使用权限\n\n你还没有获得使用权限。\n\n点击下面按钮申请使用："),
			&ext.ReplyOpts{
				Markup: &tg.ReplyInlineMarkup{
					Rows: []tg.KeyboardButtonRow{
						{
							Buttons: []tg.KeyboardButtonClass{
								&tg.KeyboardButtonCallback{
									Text: "🔐 申请使用",
									Data: []byte("request_access"),
								},
							},
						},
					},
				},
			},
		)

		return err
	}

	_, err := ctx.Reply(
		u,
		ext.ReplyTextString(`👋 欢迎使用 Shuoce File Stream

📁 发送文件、视频、图片或音频
🔗 自动生成可访问的直链
⚡ 简单、快速，无需额外操作

⚠️ 使用声明
请勿使用本机器人存储、传播任何违法违规内容。
请合理使用本服务，由此产生的一切责任由使用者自行承担。

感谢你的使用 ❤️`),
		nil,
	)

	return err
