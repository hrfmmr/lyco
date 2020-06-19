package cli

type Lyco struct {
	Help  bool `long:"help" short:"h" optional:"true" optional-value:"true" description:"Show this help message."`
	Debug bool `long:"debug" hidden:"true" optional:"true" optional-value:"true" description:"Enable lyco 🐛debugging."`
}
