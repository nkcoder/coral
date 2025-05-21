package config

import (
	"encoding/json"
	"fmt"

	"coral.daniel-guo.com/internal/aws"
	"coral.daniel-guo.com/internal/db"
)

func LoadDBConfig(env string) (*db.DBConfig, error) {
	secretName := fmt.Sprintf("hub-insights-rds-cluster-readonly-%s", env)
	secretData, err := aws.GetSecret(secretName)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret from :%s : %w", env, err)
	}

	var dbConfig db.DBConfig
	if err := json.Unmarshal([]byte(secretData), &dbConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret data: %w", err)
	}

	return &dbConfig, nil
}
