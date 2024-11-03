package msg

import (
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"sync"
)

const (
	PrivateMsg  = "private"
	GroupMsg    = "group"
	SubTextMsg  = "text"
	SubImageMsg = "image"
	SubAtMsg    = "at"
	SubReplyMsg = "reply"
)

var col *mongo.Collection
var mu sync.Mutex

type ReceiveMessage struct {
	SelfId        int64           `json:"self_id" bson:"self_id"`
	UserId        int64           `json:"user_id" bson:"user_id"`
	Time          int64           `json:"time" bson:"time"`
	MessageId     int64           `json:"message_id" bson:"message_id"`
	RealId        int64           `json:"real_id" bson:"real_id"`
	MessageSeq    int64           `json:"message_seq" bson:"message_seq"`
	MessageType   string          `json:"message_type" bson:"message_type"`
	Sender        *Sender         `json:"sender" json:"sender"`
	RawMessage    string          `json:"raw_message" bson:"raw_message"`
	Font          int64           `json:"font" bson:"font"`
	SubType       string          `json:"sub_type" bson:"sub_type"`
	Message       Messages        `json:"message" bson:"message"`
	MessageFormat string          `json:"message_format" bson:"message_format"`
	PostType      string          `json:"post_type" bson:"post_type"`
	GroupId       int64           `json:"group_id" bson:"group_id"`
	Handled       bool            `json:"-" bson:"-"`
	FunctionName  string          `json:"-" bson:"-"`
	Mutex         sync.Mutex      `json:"-" gorm:"-" bson:"-"`
	Reply         *Message        `json:"reply" bson:"-"`
	ReplyMessage  *ReceiveMessage `json:"reply_message" bson:"reply_message"`
	Entry         string          `json:"entry" gorm:"-" bson:"-"`
	Command       string          `json:"command" gorm:"-" bson:"-"`
	Text          string          `json:"text" gorm:"-" bson:"-"`
}

type Sender struct {
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card"`
}

type Message struct {
	Data map[string]interface{} `json:"data"`
	Type string                 `json:"type"`
}

type Messages []Message

func (m *ReceiveMessage) IsGroupMessage() bool {
	return m.MessageType == GroupMsg
}

func (m *ReceiveMessage) IsPrivateMessage() bool {
	return m.MessageType == PrivateMsg
}

func (m *ReceiveMessage) CheckFormat() bool {
	if len(m.Message) < 2 || m.Message[0].Type != SubAtMsg || m.Message[1].Type != SubTextMsg {
		return false
	}
	atqq, _ := strconv.ParseInt(m.Message[0].Data["qq"].(string), 10, 64)
	if m.SelfId != atqq {
		return false
	}
	return true
}

func (m *ReceiveMessage) CheckSource(whiteGroup, banGroup []int64) bool {
	if whiteGroup != nil && len(whiteGroup) > 0 {
		for _, group := range whiteGroup {
			if m.GroupId == group {
				return true
			}
		}
		return false
	}

	if banGroup != nil && len(banGroup) > 0 {
		for _, group := range banGroup {
			if m.GroupId == group {
				return false
			}
		}
	}
	return true
}

func (m *ReceiveMessage) GetMessageId() int64 {
	return m.MessageId
}

func (m *ReceiveMessage) SplitMessage() (entry string, command string, ok bool) {
	text := strings.TrimSpace(m.Message[1].Data["text"].(string))
	texts := strings.Split(text, " ")
	num := len(texts)
	if num <= 0 {
		return "", "", false
	}

	if num == 1 {
		m.Entry = texts[0]
		entry = texts[0]
		return entry, "", true
	}

	if num == 2 {
		m.Entry = texts[0]
		entry = texts[0]
		if strings.HasPrefix(texts[1], "-") {
			m.Command, _ = strings.CutPrefix(texts[1], "-")
			command, _ = strings.CutPrefix(texts[1], "-")
		} else {
			m.Text = texts[1]
		}
		return entry, command, true
	}

	if num == 3 {
		m.Entry = texts[0]
		entry = texts[0]
		if strings.HasPrefix(texts[1], "-") {
			m.Command, _ = strings.CutPrefix(texts[1], "-")
			command, _ = strings.CutPrefix(texts[1], "-")
			m.Text = texts[2]
		} else {
			m.Text = texts[1]
			m.Text = m.Text + texts[2]
		}
		return entry, command, true
	}

	return "", "", false

}
