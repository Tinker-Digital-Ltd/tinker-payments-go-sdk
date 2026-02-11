package model

type ApiMeta struct {
	RequestID   string
	Timestamp   string
	Environment string
}

func NewApiMeta(meta map[string]interface{}) *ApiMeta {
	m := &ApiMeta{}
	if meta == nil {
		return m
	}
	if v, ok := meta["request_id"].(string); ok {
		m.RequestID = v
	}
	if v, ok := meta["timestamp"].(string); ok {
		m.Timestamp = v
	}
	if v, ok := meta["environment"].(string); ok {
		m.Environment = v
	}
	return m
}
