package config

import (
	"encoding/json"
	"fmt"
)

// JobConfig ...
type JobConfig struct {
	FileProcessingConfig FileProcessingConfig `mapstructure:"file_processing"`

	FlattenConfig          SchedulerConfig `mapstructure:"flatten"`
	ExecuteTaskConfig      SchedulerConfig `mapstructure:"execute_task"`
	ExecuteGroupTaskConfig SchedulerConfig `mapstructure:"execute_group_task"`
	UpdateStatusConfig     SchedulerConfig `mapstructure:"update_status"`
}

// FileProcessingConfig ...
type FileProcessingConfig struct {
	Schedule string `mapstructure:"schedule"`
}

// SchedulerConfig ...
type SchedulerConfig struct {
	Schedule           string `mapstructure:"schedule"`
	NumDigesters       int    `mapstructure:"num_digesters"`
	NumDigestersCustom string `mapstructure:"num_digesters_custom"`
}

func (sc SchedulerConfig) GetNumDigesters(clientID int) int {
	var numDigestersCustom []NumDigestersCustom
	if err := json.Unmarshal([]byte(sc.NumDigestersCustom), &numDigestersCustom); err != nil {
		fmt.Printf("error when convert numDigestersCustom: value=%v, err=%v", sc.NumDigestersCustom, err)
		return sc.NumDigesters
	}

	for _, cfg := range numDigestersCustom {
		if cfg.ClientID == clientID {
			return cfg.Value
		}
	}

	return sc.NumDigesters
}

// NumDigestersCustom ...
type NumDigestersCustom struct {
	ClientID int `json:"client_id"`
	Value    int `json:"value"`
}
