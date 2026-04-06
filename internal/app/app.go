package app

import (
	"github.com/miloskovacevic/mergeflow/internal/config"
	"github.com/miloskovacevic/mergeflow/internal/infrastructure/gitlab"
	"github.com/miloskovacevic/mergeflow/internal/infrastructure/jira"
	"github.com/spf13/viper"
)

type App struct {
	Config config.Config
	Jira   *jira.Client
	Gitlab *gitlab.Client
}

func NewApp() (*App, error) {
	var cfg config.Config

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Config: cfg,
		Jira:   jira.NewClient(cfg.Jira.Route, cfg.Jira.Token),
		Gitlab: gitlab.NewClient(cfg.GitLab.Route, cfg.GitLab.Token),
	}, nil
}
