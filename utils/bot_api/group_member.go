package bot_api

import (
	"chatbot/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetGroupMemberList(groupId int64) ([]MemberInfo, error) {
	addr := config.Get().Server.BotAddr
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%d", addr, "/get_group_member_list?group_id=", groupId), nil)
	if err != nil {
		return nil, err
	}
	response, err := (&http.Client{Timeout: 2 * time.Second}).Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var members []MemberInfo
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, newBotResponse(&members)); err != nil {
		return nil, err
	}
	return members, nil
}

func newBotResponse(data interface{}) *BotResponse {
	res := &BotResponse{
		Status:  "",
		RetCode: 0,
		Data:    data,
	}
	return res
}

type BotResponse struct {
	Status  string      `json:"status"`
	RetCode int         `json:"retcode"`
	Data    interface{} `json:"data"`
}

type MemberInfo struct {
	GroupId         int64  `json:"group_id"`
	UserId          int64  `json:"user_id"`
	Nickname        string `json:"nickname"`
	Card            string `json:"card"`
	Sex             string `json:"sex"`
	Age             int    `json:"age"`
	Area            string `json:"area"`
	JoinTime        int64  `json:"join_time"`
	LastSentTime    int64  `json:"last_sent_time"`
	Level           string `json:"level"`
	Role            string `json:"role"`
	Unfriendly      bool   `json:"unfriendly"`
	Title           string `json:"title"`
	TitleExpireTime int    `json:"title_expire_time"`
	CardChangeable  bool   `json:"card_changeable"`
}
