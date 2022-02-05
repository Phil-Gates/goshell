package main


import (
    "os"
    "os/exec"
    "log"
    "fmt"
    "bufio"
    "runtime"
    "errors"
)


var session_path string
const version string = "v0.0.1"


func main() {
    if runtime.GOOS == "linux" {
        session_path = "/tmp/goshell_session.go"
    } else if runtime.GOOS == "darwin" {
        session_path = "~/goshellsessions/goshell_session.go"
    } else if runtime.GOOS == "windows" {
        session_path = "$env:USERPROFILE\\goshellsessions\\goshell_session.go"
    } else {
        log.Fatal(errors.New("Unkown OS..."))
    }
    var ui string
    reader := bufio.NewReader(os.Stdin)
    init_file()
    fmt.Println("GoShell " + version)
    fmt.Println("Type '!!help' for more info")
    for {
        fmt.Print(">>> ")
        ui, _ = reader.ReadString('\n')

        if ui == "!!help\n" {
            fmt.Println("!!help: show this help message")
            fmt.Println("!!exit: exit the shell")
            fmt.Println("!!run: run code")
            fmt.Println("!!runc: run then clear buffer")
            fmt.Println("!!clrbuf: clear the buffer")
            fmt.Println("!!cache: cache the current code")
            fmt.Println("!!restore: run the cached code")
        } else if ui == "!!exit\n" {
            clr_buf()
            break
        } else if ui == "!!run\n" {
            append_session("}")
            run_go()
            del_last()
        } else if ui == "!!runc\n" {
            append_session("}")
            run_go()
            init_file()
        } else if ui == "!!clrbuf\n" {
            init_file()
        } else if ui == "!!cache\n" {
            cache_session()
        } else if ui == "!!restore\n" {
            restore_session()
        } else {
            append_session(ui)
        }
    }
}


func check_errs(err error) {
    if err != nil {
        log.Fatal(err)
    }
}


func append_session(str string) {
    session, err := os.OpenFile(session_path, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
    check_errs(err)
    _, err = session.Write([]byte(str))
    check_errs(err)
    session.Close()
}


func cache_session() {
    if _, err := os.Stat(session_path + ".gz"); os.IsNotExist(err) {
        cmd := exec.Command("gzip", "-k", session_path)
        check_errs(cmd.Run())
    } else {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Override previous cache? [y/n] ")
        confirm, _ := reader.ReadString('\n')
        if confirm == "y\n" || confirm == "Y\n" {
            os.Remove(session_path + ".gz")
            cache_session()
        } else {
            fmt.Println("Abort")
        }
    }
}


func restore_session() {
    cmd := exec.Command("gzip", "-fdk", session_path + ".gz")
    check_errs(cmd.Run())
}


func run_go() {
    cmd := exec.Command("go", "run", session_path)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    check_errs(cmd.Run())
}


func clr_buf() {
    if _, err := os.Stat(session_path); !os.IsNotExist(err) {
        os.Remove(session_path)
    }
}


func imports() {
    if len(os.Args) > 2 {
        if os.Args[1] == "--imports" {
            for i := 0; i < len(os.Args[2:]); i++ {
                append_session("import \"" + os.Args[2 + i] + "\"\n")
            }
        }
    }
}


func del_last() {
    session, err := os.Open(session_path)
    check_errs(err)
    defer session.Close()

    var lines []string
    scanner := bufio.NewScanner(session)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    check_errs(scanner.Err())

    session, err = os.OpenFile(session_path, os.O_WRONLY, 0644)
    check_errs(err)
    err = session.Truncate(0)
    check_errs(err)
    _, err = session.Seek(0, 0)
    check_errs(err)
    writer := bufio.NewWriter(session)
    for _, line := range lines[0:len(lines)-1] {
        fmt.Fprintln(writer, line)
    }
    check_errs(writer.Flush())
}


func init_file() {
    clr_buf()
    append_session("package main\n")
    imports()
    append_session("func main() {\n")
}
