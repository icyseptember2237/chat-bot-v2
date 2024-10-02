package msg

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type SendMessage interface {
	AppendSegment(seg ...*Segment)
	Send(url, token string) (*SendResult, error)
}

type Segment struct {
	Type string                 `json:"type" mapstructure:"Type"`
	Data map[string]interface{} `json:"data" mapstructure:"Data"`
}

type GroupMessage struct {
	GroupId    int64      `json:"group_id" mapstructure:"GroupId"`
	Message    []*Segment `json:"message" mapstructure:"Message"`
	AutoEscape bool       `json:"auto_escape" mapstructure:"AutoEscape"`
}

type PrivateMessage struct {
	UserId     int64      `json:"user_id" mapstructure:"UserId"`
	Message    []*Segment `json:"message" mapstructure:"Message"`
	AutoEscape bool       `json:"auto_escape" mapstructure:"AutoEscape"`
}

type SendResult struct {
	Status  string                 `json:"status"`
	Retcode int                    `json:"retcode"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func NewTextSegment(text string) *Segment {
	return &Segment{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	}
}

func NewAtSegment(at int64) *Segment {
	return &Segment{
		Type: "at",
		Data: map[string]interface{}{
			"qq": at,
		},
	}
}

func NewReplySegment(msgId int64) *Segment {
	return &Segment{
		Type: "reply",
		Data: map[string]interface{}{
			"id": msgId,
		},
	}
}

func NewImageSegment(image string) *Segment {
	return &Segment{
		Type: "image",
		Data: map[string]interface{}{
			"file": image,
		},
	}
}

func NewFileSegment(file, name string) *Segment {
	return &Segment{
		Type: "file",
		Data: map[string]interface{}{
			"file": file,
			"name": name,
		},
	}
}

func NewQQMusicSegment(id string) *Segment {
	return &Segment{
		Type: "music",
		Data: map[string]interface{}{
			"type": "qq",
			"id":   id,
		},
	}
}

func NewRecordSegment(file string) *Segment {
	return &Segment{
		Type: "record",
		Data: map[string]interface{}{
			"file": file,
		},
	}
}

func NewNetEaseMusicSegment(id string) *Segment {
	return &Segment{
		Type: "music",
		Data: map[string]interface{}{
			"type": "163",
			"id":   id,
		},
	}
}

func NewGroupMessage(receiver int64, segments ...*Segment) *GroupMessage {
	msg := &GroupMessage{
		GroupId: receiver,
		Message: make([]*Segment, 0),
	}
	for _, segment := range segments {
		msg.Message = append(msg.Message, segment)
	}
	return msg
}

func (msg *GroupMessage) AppendSegment(seg ...*Segment) {
	msg.Message = append(msg.Message, seg...)
}

func (msg *GroupMessage) Send(url, token string) (*SendResult, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url+"/send_group_msg", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	res, err := (&http.Client{Timeout: 3 * time.Second}).Do(request)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var body SendResult
	if err := json.Unmarshal(resBody, &body); err != nil {
		return nil, err
	}
	return &body, nil
}

func NewPrivateMessage(receiver int64, segments ...*Segment) *PrivateMessage {
	msg := &PrivateMessage{
		UserId:  receiver,
		Message: make([]*Segment, 0),
	}
	for _, segment := range segments {
		msg.Message = append(msg.Message, segment)
	}
	return msg
}

func (msg *PrivateMessage) AppendSegment(seg ...*Segment) {
	msg.Message = append(msg.Message, seg...)
}

func (msg *PrivateMessage) Send(url, token string) (*SendResult, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url+"/send_private_msg", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	res, err := (&http.Client{Timeout: 3 * time.Second}).Do(request)
	if err != nil {
		return nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var body SendResult
	if err := json.Unmarshal(resBody, &body); err != nil {
		return nil, err
	}
	return &body, nil
}
