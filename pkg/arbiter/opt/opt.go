package opt

type Option struct {
	TriggerKeepAlive int `json:"trigger_keepalive"`
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) SetTriggerKeepalive(triggerKeepalive int) {
	o.TriggerKeepAlive = triggerKeepalive
}

func (o *Option) GetTriggerKeepalive() int {
	return o.TriggerKeepAlive
}
