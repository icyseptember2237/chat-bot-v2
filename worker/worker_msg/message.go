package worker_msg

type Message struct {
	Info    Info    `json:"title" mapstructure:"title"`
	Content Content `json:"content" mapstructure:"content"`
}

type Info struct {
	From        string `json:"from" mapstructure:"from"`
	To          string `json:"to" mapstructure:"to"`
	Time        int64  `json:"time" mapstructure:"time"`
	ProcessTime int64  `json:"process_time" mapstructure:"process_time"`
	FinishTime  int64  `json:"finish_time" mapstructure:"finish_time"`
}

type Content struct {
	Data map[string]interface{} `json:"data" mapstructure:"data"`
}
