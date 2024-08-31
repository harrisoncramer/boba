package boba

type KeyOpts struct {
	Up     string `mapstructure:"up"`
	Down   string `mapstructure:"down"`
	Select string `mapstructure:"select"`
	Toggle string `mapstructure:"toggle"`
	Back   string `mapstructure:"back"`
	Quit   string `mapstructure:"quit"`
	Filter string `mapstructure:"filter"`
	Help   string `mapstructure:"help"`
}
