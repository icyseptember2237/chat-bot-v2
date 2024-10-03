package config

type WorkerConfig struct {
	Name      string `json:"name" mapstructure:"name"`
	Enable    bool   `json:"enable" mapstructure:"enable"`
	Num       int    `json:"num" mapstructure:"num"`
	RateLimit int    `json:"rate_limit" mapstructure:"rate_limit"`
	Script    string `json:"script" mapstructure:"script"`
	Handler   string `json:"handler" mapstructure:"handler"`

	Source *Source `json:"source" mapstructure:"source"`
	Dest   *Dest   `json:"dest" mapstructure:"dest"`

	Config map[string]interface{} `json:"config" mapstructure:"config"`
}

type Source struct {
	Type         string                 `json:"type" mapstructure:"type"`
	BufferLength int                    `json:"buffer_length" mapstructure:"buffer_length"`
	Config       map[string]interface{} `json:"config" mapstructure:"config"`
}

type Dest struct {
	Type         string                 `json:"type" mapstructure:"type"`
	BufferLength int                    `json:"buffer_length" mapstructure:"buffer_length"`
	Config       map[string]interface{} `json:"config" mapstructure:"config"`
}
