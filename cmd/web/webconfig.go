package main

type webconfig struct {
	WebUsers map[string]string
	ApiUsers map[string]string

	ConfigFileName      string
	PersonalPinFileName string
	AccessLogFileName   string
	ListenPort          string
}
