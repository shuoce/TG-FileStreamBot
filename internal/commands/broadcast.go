package commands

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/authorization"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/tg"
)

var (
	broadcastWaiting   bool
	broadcastWaitingMu sync.Mutex
)

func (m *command) LoadBroadcast(d dispatcher.Dispatcher) {
	log := m.log.Named("broadcast")
	defer log.Sugar().Info("Loaded")

	// /broadcast
	d.AddHandler(
		handlers.NewCommand("broadcast", broadcast),
	)

	// /cancel
	d.AddHandler(
		handlers.NewCommand("cancel", cancelBroadcast),
	)

	// 广播内容
	d.AddHandler(
		handlers.NewMessage(nil, handleBroadcastMessage),
	)
}

func broadcast(ctx *ext.Context, u *ext.Update) error {
	if u.EffectiveChat().GetID() != config.ValueOf.AdminID {
		return nil
	}

	broadcastWaitingMu.Lock()
	broadcastWaiting = true
	broadcastWaitingMu.Unlock()

	_, err := ctx.Reply(
		u,
		ext.ReplyTextString(
			"📢 广播模式\n\n请输入要发送给所有已授权用户的通知内容。\n\n发送 /cancel 可取消广播。",
		),
		nil,
	)

	return err
}

func cancelBroadcast(ctx *ext.Context, u *ext.Update) error {
	if u.EffectiveChat().GetID() != config.ValueOf.AdminID {
		return nil
	}

	broadcastWaitingMu.Lock()
	broadcastWaiting = false
	broadcastWaitingMu.Unlock()

	_, err := ctx.Reply(
		u,
		ext.ReplyTextString("✅ 已取消广播。"),
		nil,
	)

	if err != nil {
		return err
	}

	return dispatcher.EndGroups
}

func handleBroadcastMessage(ctx *ext.Context, u *ext.Update) error {
	if u.EffectiveChat().GetID() != config.ValueOf.AdminID {
		return nil
	}

	broadcastWaitingMu.Lock()
	waiting := broadcastWaiting

	if waiting {
		broadcastWaiting = false
	}

	broadcastWaitingMu.Unlock()

	if !waiting {
		return nil
	}

	message := ""

if msg := u.EffectiveMessage; msg != nil {
    message = strings.TrimSpace(msg.Message)
}

	if message == "" {
		broadcastWaitingMu.Lock()
		broadcastWaiting = true
		broadcastWaitingMu.Unlock()

		_, err := ctx.Reply(
			u,
			ext.ReplyTextString("❌ 广播内容不能为空，请重新发送内容。\n\n发送 /cancel 可取消。"),
			nil,
		)

		return err
	}

	users, err := authorization.GetApprovedUsers()
	if err != nil {
		_, replyErr := ctx.Reply(
			u,
			ext.ReplyTextString("❌ 获取用户列表失败："+err.Error()),
			nil,
		)
		return replyErr
	}

	if len(users) == 0 {
		_, err := ctx.Reply(
			u,
			ext.ReplyTextString("📭 当前没有已授权用户。"),
			nil,
		)
		return err
	}

	_, err = ctx.Reply(
		u,
		ext.ReplyTextString(
			fmt.Sprintf(
				"📢 广播开始\n\n👥 用户数：%d\n⏳ 正在发送，请稍候……",
				len(users),
			),
		),
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		success := 0
		failed := 0

		for _, userID := range users {
			_, sendErr := ctx.SendMessage(
				userID,
				&tg.MessagesSendMessageRequest{
					Message: message,
				},
			)

			if sendErr != nil {
				failed++
			} else {
				success++
			}

			// 控制发送速度
			time.Sleep(40 * time.Millisecond)
		}

		result := fmt.Sprintf(
			"📢 广播完成\n\n👥 总人数：%d\n✅ 成功：%d\n❌ 失败：%d",
			len(users),
			success,
			failed,
		)

		_, _ = ctx.SendMessage(
			config.ValueOf.AdminID,
			&tg.MessagesSendMessageRequest{
				Message: result,
			},
		)
	}()

	return dispatcher.EndGroups
}
