package bot_api

import (
	"bytes"
	"chatbot/config"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type botResp[D any] struct {
	Status  string `json:"status"`
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    *D     `json:"data"`
}

func RequestBot[D any](method string, api string, req any, respData *D) error {
	conf := config.Get().Server
	reqData, err := json.Marshal(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(method, conf.BotAddr+api, bytes.NewBuffer(reqData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+conf.BotToken)
	res, err := (&http.Client{Timeout: 3 * time.Second}).Do(request)
	if err != nil {
		return err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	resp := botResp[D]{
		Data: respData,
	}
	if err = json.Unmarshal(resBody, &resp); err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New(resp.Message)
	}
	return nil
}
