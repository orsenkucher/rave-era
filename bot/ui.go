package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/rave-era/repo"
)

func genFreeList(list []repo.Event) (tgbotapi.InlineKeyboardMarkup, bool) {
	buttons := make([]tgbotapi.InlineKeyboardButton, len(list))
	for i, e := range list {
		icon := "🗽"
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf(" %s %s  ", icon, e.Name), "free:"+e.Name)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, b := range buttons {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{b})
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...), true
}
