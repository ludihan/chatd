package cli

import (
    "chatd/consts"
    "fmt"
    "os"
)

func HandleArgs() {
    if len(os.Args) < 2 {
        fmt.Println(consts.Help)
        os.Exit(0)
    } else {
        switch os.Args[1] {
        case "-h", "--help":
            fmt.Println(consts.Help)
            os.Exit(0)
        case "-v", "--version":
            fmt.Println(consts.Version)
            os.Exit(0)
        case "-d", "--default":
            if len(os.Args) < 3 {
                fmt.Println("error: argument to -d missing target file")
            } else {

            }

        }
    }
}
