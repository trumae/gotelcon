package engine

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/vvmaciel/mundocrypto/db"
)

//EntryChatContext ...
type EntryChatContext struct {
	Text    string
	Regexp  string //For future use
	Handler func(req *tgbotapi.Message, msg *tgbotapi.MessageConfig, user *db.User)
	Ubi     bool
}

//ChatContext ...
type ChatContext struct {
	Hs             []EntryChatContext
	FnNoCommand    func(req *tgbotapi.Message, msg *tgbotapi.MessageConfig, user *db.User)
	State          int
	TransferTo     string
	TransferAmount float64
}

//Clear ..
func (ec *ChatContext) Clear() {
	n := ChatContext{}
	for _, v := range ec.Hs {
		if v.Ubi {
			n.Hs = append(n.Hs, v)
		}
	}
	ec.Hs = n.Hs
}

//Add ...
func (ec *ChatContext) Add(txt string, f func(req *tgbotapi.Message, msg *tgbotapi.MessageConfig, user *db.User)) {
	n := EntryChatContext{
		Text:    txt,
		Handler: f,
		Ubi:     false,
	}

	ec.Hs = append(ec.Hs, n)
}

//AddUbi ..
func (ec *ChatContext) AddUbi(txt string, f func(req *tgbotapi.Message, msg *tgbotapi.MessageConfig, user *db.User)) {
	n := EntryChatContext{
		Text:    txt,
		Handler: f,
		Ubi:     true,
	}

	ec.Hs = append(ec.Hs, n)
}

//Contexts ...
var Contexts = make(map[int]*ChatContext)

//GetOrNewChatContext ...
func GetOrNewChatContext(id int, base *ChatContext) *ChatContext {
	c, ok := Contexts[id]
	if !ok {
		n := &ChatContext{}
		for _, v := range base.Hs {
			n.Hs = append(n.Hs, v)
		}
		Contexts[id] = n
		return n
	}

	return c
}

//GetChatContext ...
func GetChatContext(req *tgbotapi.Message) (*ChatContext, error) {
	id := req.From.ID
	c, ok := Contexts[id]
	if !ok {
		return nil, fmt.Errorf("id not found")
	}

	return c, nil
}

//HandleMessages ...
func HandleMessages(req *tgbotapi.Message, msg *tgbotapi.MessageConfig, base *ChatContext, user *db.User) {
	txt := req.Text
	//log.Println("--- txt", txt)

	chat := GetOrNewChatContext(req.From.ID, base)
	//log.Println("--- chat", chat)

	for _, v := range chat.Hs {
		//log.Println("--- v.Text", v.Text)
		if v.Text == txt {
			//log.Println("--- found")
			v.Handler(req, msg, user)
			return
		}
	}
	base.FnNoCommand(req, msg, user)
}
