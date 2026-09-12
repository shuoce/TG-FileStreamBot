package commands

import (
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
		ext.ReplyTextString("你好，欢迎使用此机器人，给我发一个 文件/视频/图片/音频 我可以帮你生成直链。"),
		nil,
	)

	return err
}
