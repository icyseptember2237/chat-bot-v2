package msg

import (
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

type ReceiveMessage struct {
	SelfId        int64      `json:"self_id"`
	UserId        int64      `json:"user_id"`
	Time          int64      `json:"time"`
	MessageId     int64      `json:"message_id"`
	RealId        int64      `json:"real_id"`
	MessageSeq    int64      `json:"message_seq"`
	MessageType   string     `json:"message_type"`
	Sender        *Sender    `json:"sender" gorm:"type:jsonb"`
	RawMessage    string     `json:"raw_message"`
	Font          int64      `json:"font"`
	SubType       string     `json:"sub_type"`
	Message       Messages   `json:"message" gorm:"type:jsonb"`
	MessageFormat string     `json:"message_format"`
	PostType      string     `json:"post_type"`
	GroupId       int64      `json:"group_id"`
	Handled       bool       `gorm:"type:boolean"`
	FunctionName  string     `json:"function_name"`
	Mutex         sync.Mutex `gorm:"-" json:"-"`
	Entry         string     `gorm:"-" json:"entry"`
	Command       string     `gorm:"-" json:"command"`
	Text          string     `gorm:"-" json:"text"`
}

type Sender struct {
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card"`
}

type Message struct {
	Data map[string]string `json:"data"`
	Type string            `json:"type"`
}

type Messages []Message

func (m *ReceiveMessage) IsGroupMessage() bool {
	return m.MessageType == GroupMsg
}

func (m *ReceiveMessage) IsPrivateMessage() bool {
	return m.MessageType == PrivateMsg
}

func (m *ReceiveMessage) CheckFormat() bool {
	if !m.IsGroupMessage() {
		return false
	}
	if len(m.Message) < 2 || m.Message[0].Type != SubAtMsg || m.Message[1].Type != SubTextMsg {
		return false
	}
	atqq, _ := strconv.ParseInt(m.Message[0].Data["qq"], 10, 64)
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
	text := strings.TrimSpace(m.Message[1].Data["text"])
	texts := strings.Split(text, " ")
	num := len(texts)
	if num <= 0 {
		return "", "", false
	}
	if !strings.HasPrefix(texts[0], "/") {
		m.Entry = "/help"
		entry = "/help"
		m.Command = ""
		m.Text = text
		return entry, command, true
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
