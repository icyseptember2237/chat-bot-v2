package msg

import (
	my_mongo "chatbot/storage/mongo"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"sync"
	"time"
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
	ReplyMessage  *ReceiveMessage `json:"reply_message" bson:"-"`
	Entry         string          `json:"entry" gorm:"-" bson:"-"`
	Command       string          `json:"command" gorm:"-" bson:"-"`
	ArgAt         []int64         `json:"arg_at" gorm:"-" bson:"-"`
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

func (m *ReceiveMessage) GetMessageId() int64 {
	return m.MessageId
}

func (m *ReceiveMessage) SplitMessage() (entry string, command string, ok bool) {
	// 将消息段中的作为参数的at取出
	// example: @bot /func -command @user
	if len(m.Message) == 3 && m.Message[2].Type == SubAtMsg {
		at := m.Message[3].Data["qq"].(string)
		qq, _ := strconv.ParseInt(at, 10, 64)
		m.ArgAt = append(m.ArgAt, qq)
	}

	// 处理消息段
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

func (m *ReceiveMessage) ResolveMessage() bool {
	segments := m.Message
	atBot, functioned, commanded := false, false, false

	for _, segment := range segments {
		switch segment.Type {
		case SubReplyMsg:
			m.Reply = &Message{
				Data: segment.Data,
				Type: segment.Type,
			}
			if col == nil && mu.TryLock() {
				defer mu.Unlock()
				col = my_mongo.Get("default").Database("bot").Collection("msg")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if msgId, err := strconv.ParseInt(segment.Data["id"].(string), 10, 64); err == nil {
				if err := col.FindOne(ctx, bson.M{"message_id": msgId}).Decode(&m.ReplyMessage); err != nil {
					fmt.Println(err)
				}
			}
			break
		case SubAtMsg:
			if !atBot {
				atqq, _ := strconv.ParseInt(segment.Data["qq"].(string), 10, 64)
				if m.SelfId == atqq {
					atBot = true
				}
			} else if functioned {
				at := segment.Data["qq"].(string)
				qq, _ := strconv.ParseInt(at, 10, 64)
				m.ArgAt = append(m.ArgAt, qq)
			}
			break
		case SubTextMsg:
			if atBot {
				if !functioned {
					text := strings.TrimSpace(segment.Data["text"].(string))
					for _, v := range strings.Split(text, " ") {
						if strings.HasPrefix(v, "/") {
							m.Entry = v
							functioned = true
						} else if functioned {
							if strings.HasPrefix(v, "-") && !commanded {
								m.Command, _ = strings.CutPrefix(v, "-")
								commanded = true
							} else {
								m.Text = v
							}
						}
					}
				} else {
					m.Text = m.Text + segment.Data["text"].(string)
				}
			}
			break
		default:
			break
		}
	}
	return atBot && functioned
}
