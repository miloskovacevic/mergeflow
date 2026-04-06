package config

type Config struct {
	Jira struct {
		Route   string
		Project string
		Token   string
	}

	GitLab struct {
		Route string
		Repos []int
		Token string
	}
}
