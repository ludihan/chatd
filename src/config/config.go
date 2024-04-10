package config

type Session struct {
    Visible bool
    User_list bool
    User_max bool
    Filter_token []string
    Filter_whole []string
}

type Daemon struct {
    Default Session
    Sessions map[string]Session
}
