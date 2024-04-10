package consts

var Help =
`Usage: chatd [OPTIONS] FILE

Options:
    -h              Display this message
    -v              Show version info
    -d CONFIGFILE   Creates a default toml config file as CONFIGFILE

Examples:
    chatd -d ./config.toml
    chatd ./config.toml`

var Version =
`chatd version 0.0.1`
