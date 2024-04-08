package config

type Session struct {

}

type Daemon struct {
    Default Session
    Sessions map[string]Session
}

const (
    FileNotFound = iota
    SyntaxError
    SemanticError
)
