package config

type JobConfig struct {
	Name       string  `json:"name" mapstructure:"name"`
	Enable     bool    `json:"enable" mapstructure:"enable"`
	WhiteGroup []int64 `json:"white_group" mapstructure:"white_group"`
	Cron       string  `json:"cron" mapstructure:"cron"`
	Script     string  `json:"script" mapstructure:"script"`
	Handler    string  `json:"handler" mapstructure:"handler"`

	Config map[string]interface{} `json:"config" mapstructure:"config"`
}
