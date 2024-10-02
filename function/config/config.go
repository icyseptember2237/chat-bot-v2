package config

import (
	engine_pool "chatbot/utils/engine_pool"
	"golang.org/x/time/rate"
)

type Handler struct {
	Command     string `mapstructure:"command" json:"command"`
	Description string `mapstructure:"description" json:"description"`
	Script      string `mapstructure:"script" json:"script"`
	Handler     string `mapstructure:"handler" json:"handler"`
	RateLimit   int    `mapstructure:"rate_limit" json:"rate_limit"`

	rl *rate.Limiter           `mapstructure:"-" json:"-"`
	ep *engine_pool.EnginePool `mapstructure:"-" json:"-"`
}

func (h *Handler) GetEnginePool() *engine_pool.EnginePool {
	return h.ep
}

func (h *Handler) SetEnginePool(ep *engine_pool.EnginePool) {
	h.ep = ep
}

func (h *Handler) GetLimiter() *rate.Limiter {
	return h.rl
}

func (h *Handler) SetLimiter(rl *rate.Limiter) {
	h.rl = rl
}

func (h *Handler) ReachLimit() bool {
	if h.rl == nil {
		return false
	}
	return !h.rl.Allow()
}

type Function struct {
	Name        string    `mapstructure:"name" json:"name"`
	Description string    `mapstructure:"description" json:"description"`
	Script      string    `mapstructure:"script" json:"script"`
	Handlers    []Handler `mapstructure:"handlers" json:"handlers"`
}

type Static struct {
	Path string `mapstructure:"path" json:"path"`
	Root string `mapstructure:"root" json:"root"`
}

type Server struct {
	ServerPort  string  `mapstructure:"server_port" json:"server_port"`
	ServerAddr  string  `mapstructure:"server_addr" json:"server_addr"`
	ServerToken string  `mapstructure:"server_token" json:"server_token"`
	BotAddr     string  `mapstructure:"bot_addr" json:"bot_addr"`
	BotToken    string  `mapstructure:"bot_token" json:"bot_token"`
	BotNumber   string  `mapstructure:"bot_number" json:"bot_number"`
	AdminPort   string  `mapstructure:"admin_port" json:"admin_port"`
	Static      *Static `mapstructure:"static" json:"static"`

	PreloadCnt int `mapstructure:"preload_cnt" json:"preload_cnt"`
	// 全局设置字段
	Config     map[string]interface{} `mapstructure:"config" json:"config"`
	WhiteGroup []int64                `mapstructure:"white_group" json:"white_group"`
	BanGroup   []int64                `mapstructure:"ban_group" json:"ban_group"`
	Functions  []Function             `mapstructure:"functions" json:"functions"`
}
