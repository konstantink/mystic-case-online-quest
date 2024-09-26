package main

import (
	"fmt"
	"log"
	"mystic-case-online-quest/config"
	"strconv"
	"strings"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
)

type BotMessage struct {
	UserID uuid.UUID
	Page   string
	IsNew  bool
}

func (oq *OnlineQuest) initBot() {
	var err error
	pref := tele.Settings{
		Token: config.Config("BOT_TOKEN"),
	}

	oq.bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatal("[ERROR] failed to init bot")
		return
	}

	oq.botChan = make(chan BotMessage, 1)
}

func (oq *OnlineQuest) ListenForMessages() {
	for message := range oq.botChan {
		log.Print("[INFO] sending message to the channel")
		chatId, _ := strconv.ParseInt(config.Config("BOT_CHAT_ID"), 10, 64)
		chat, _ := oq.bot.ChatByID(chatId)
		oq.bot.Send(chat, fmt.Sprintf("*User:* _#%s_\n*Is new:* %t\n*Page:* %s\n", strings.Replace(message.UserID.String(), "-", "", -1), message.IsNew, message.Page), &tele.SendOptions{ParseMode: tele.ModeMarkdown})
	}
}
