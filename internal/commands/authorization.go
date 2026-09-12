package commands

import (
	"fmt"
	"strconv"
	"strings"

	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/authorization"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/tg"
)

func (m *command) LoadAuthorization(dispatcher dispatcher.Dispatcher) {
	log := m.log.Named("authorization")
	defer log.Sugar().Info("Loaded")

	dispatcher.AddHandler(
		handlers.NewCallbackQuery(
			filters.CallbackQuery.Equal("request_access"),
			requestAccess,
		),
	)

	dispatcher.AddHandler(
		handlers.NewCallbackQuery(
			filters.CallbackQuery.Prefix("auth_approve:"),
			approveUser,
		),
	)

	dispatcher.AddHandler(
		handlers.NewCallbackQuery(
			filters.CallbackQuery.Prefix("auth_reject:"),
			rejectUser,
		),
	)
}

func requestAccess(ctx *ext.Context, u *ext.Update) error {
	query := u.CallbackQuery
	userID := query.UserID

	_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: query.QueryID,
		Message: "申请已提交，请等待管理员审核。",
		Alert:   true,
	})

	username := ""
	if user := u.EffectiveUser(); user != nil {
		username = user.Username
	}

	userInfo := fmt.Sprintf(
		"🔔 新的使用权限申请\n\n👤 用户：%d\n🔗 用户名：@%s\n🆔 ID：%d",
		userID,
		username,
		userID,
	)

	_, err := ctx.SendMessage(config.ValueOf.AdminID, &tg.MessagesSendMessageRequest{
		Message: userInfo,
		ReplyMarkup: &tg.ReplyInlineMarkup{
			Rows: []tg.KeyboardButtonRow{
				{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButtonCallback{
							Text: "✅ 批准",
							Data: []byte("auth_approve:" + strconv.FormatInt(userID, 10)),
						},
						&tg.KeyboardButtonCallback{
							Text: "❌ 拒绝",
							Data: []byte("auth_reject:" + strconv.FormatInt(userID, 10)),
						},
					},
				},
			},
		},
	})

	return err
}

func approveUser(ctx *ext.Context, u *ext.Update) error {
	return handleApproval(ctx, u, true)
}

func rejectUser(ctx *ext.Context, u *ext.Update) error {
	return handleApproval(ctx, u, false)
}

func handleApproval(ctx *ext.Context, u *ext.Update, approve bool) error {
	query := u.CallbackQuery

	if query.UserID != config.ValueOf.AdminID {
		_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
			QueryID: query.QueryID,
			Message: "你没有管理员权限。",
			Alert:   true,
		})
		return nil
	}

	data := string(query.Data)
	parts := strings.SplitN(data, ":", 2)
	if len(parts) != 2 {
		return nil
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil
	}

	if approve {
		if err := authorization.Approve(userID); err != nil {
			_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
				QueryID: query.QueryID,
				Message: "批准失败。",
				Alert:   true,
			})
			return err
		}

		_, _ = ctx.SendMessage(userID, &tg.MessagesSendMessageRequest{
			Message: "🎉 审核通过！\n\n你现在可以使用文件直链机器人了。",
		})

		_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
			QueryID: query.QueryID,
			Message: "已批准该用户。",
			Alert:   true,
		})
	} else {
		_, _ = ctx.SendMessage(userID, &tg.MessagesSendMessageRequest{
			Message: "❌ 申请未通过。",
		})

		_, _ = ctx.AnswerCallback(&tg.MessagesSetBotCallbackAnswerRequest{
			QueryID: query.QueryID,
			Message: "已拒绝该用户。",
			Alert:   true,
		})
	}

	return nil
}
