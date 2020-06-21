package cli

type Lyco struct {
	PomodoroDuration    string `long:"duration" short:"d" description:"🍅 Pomdoro timer's duration (eg. 25m)"`
	ShortBreaksDuration string `long:"short-breaks" description:"☕ Short breaks duration (eg. 5m)"`
	LongBreaksDuration  string `long:"long-breaks" description:"☕ Long breaks duration (eg. 15m)"`
	Help                bool   `long:"help" short:"h" optional:"true" optional-value:"true" description:"Show this help message."`
	Debug               bool   `long:"debug" hidden:"true" optional:"true" optional-value:"true" description:"Enable lyco 🐛debugging."`
}
