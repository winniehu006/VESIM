package notify

import "github.com/winniehu006/vesim/internal/config"

type FeishuChannelOptions struct {
	Config     config.FeishuConfig
	StateStore RuntimeStateStore
}
