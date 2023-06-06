package types

type Option struct {
	Host       string
	Target     string
	ConfigFile string
}

type ImageList struct {
	images []string `yaml:"images"`
}
