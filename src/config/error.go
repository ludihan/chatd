package config

type configError int

const (
    FileNotFound = iota
    SyntaxError
    MissingParameters
)

func (c configError) Error() string {
    switch c {
    case FileNotFound:
    case SyntaxError:
    case MissingParameters:
    }
}
