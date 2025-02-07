package hera

type Config struct {
	Services map[string]*Service `yaml:"services"`
}

type Service struct {
	Command string   `yaml:"command"`
	Watch   []string `yaml:"watch"`
	Exclude []string `yaml:"exclude"`
}
