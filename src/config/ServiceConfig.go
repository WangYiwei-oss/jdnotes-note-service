package config

import "github.com/WangYiwei-oss/jdnotes-note-service/src/services"

type MServiceConfig struct {
}

func NewMServiceConfig() *MServiceConfig {
	return &MServiceConfig{}
}

func (s *MServiceConfig) JdInitCommonService() *services.CommonService {
	return services.NewCommonService()
}

func (s *MServiceConfig) JdInitNotifyProcessor() *services.NotifyProcessor {
	return services.NewNotifyProcessor()
}
