package webhook

type Meta struct {
	Version string
	AppID   string
	Gateway *string
}

func NewMeta(meta map[string]interface{}) *Meta {
	m := &Meta{
		Version: "1.0",
		AppID:   "",
	}

	if version, ok := meta["version"].(string); ok {
		m.Version = version
	}

	if appID, ok := meta["app_id"].(string); ok {
		m.AppID = appID
	}

	if gateway, ok := meta["gateway"].(string); ok {
		m.Gateway = &gateway
	}

	return m
}
