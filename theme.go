package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "text/template"
)

type PywalJsonSpecial struct {
    Foreground string `json:"foreground"`
    Background string `json:"background"`
    Cursor string `json:"cursor"`
}

type PywalJsonColors struct {
    Color0 string `json:"color0"`
    Color1 string `json:"color1"`
    Color2 string `json:"color2"`
    Color3 string `json:"color3"`
    Color4 string `json:"color4"`
    Color5 string `json:"color5"`
    Color6 string `json:"color6"`
    Color7 string `json:"color7"`
    Color8 string `json:"color8"`
    Color9 string `json:"color9"`
    Color10 string `json:"color10"`
    Color11 string `json:"color11"`
    Color12 string `json:"color12"`
    Color13 string `json:"color13"`
    Color14 string `json:"color14"`
    Color15 string `json:"color15"`
}

// Struct to demarshal (decode) json themes from pywal repo
type PywalJsonTheme struct {
    Special PywalJsonSpecial `json:"special"`
    Colors PywalJsonColors `json:"colors"`
}

// XResources theme
type XRTheme struct {
    Entries map[string]string
}

const XRDB_FILENAME = "Xtheme"

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    xdgConfig, err := os.UserConfigDir()
    check(err)

    if (os.Args[1] == "help") {
        print_help()
    }

    if (len(os.Args) < 2) {
        fmt.Println("Invalid usage: Specify a command.")
        os.Exit(1)
    }

    if (os.Args[1] == "set") {
        if (len(os.Args) != 3) {
            fmt.Println("Usage: theme set filename")
            os.Exit(1)
        }

        filename := os.Args[2]
        var theme = load(filename)
        set(theme)
    } else if (os.Args[1] == "save") {
        if (len(os.Args) != 3) {
            fmt.Println("Usage: theme save filename")
            os.Exit(1)
        }

        filename := os.Args[2]
        var theme = load(xdgConfig + "/" + XRDB_FILENAME)
        for id, value := range theme.Entries {
            fmt.Println(id, value)
        }

        save(theme, filename)
    } else {
        fmt.Println("Usage: theme set filename")
        os.Exit(1)
    }
}

func print_help() {
    fmt.Println("Available commands: help, save, set")
    fmt.Println("   Save current theme to file: theme save filename")
    fmt.Println("   Load and set current theme from file: theme set filename")
    fmt.Println("Example usage:")
}

func load(filename string) XRTheme {
    if _, err := os.Stat(filename); err != nil {
        fmt.Println("Error reading file. Does the file exist?")
        os.Exit(1)
    }

    filedata, err := ioutil.ReadFile(filename)
    check(err)

    var theme XRTheme
    suffix := fileSuffix(filename)
    if (suffix == "json") {
        // Assumed pywal-type of json
        var pwTheme PywalJsonTheme
        err = json.Unmarshal(filedata, &pwTheme)
        check(err)

        theme = pwTheme.toXRTheme()
    } else if (suffix == "") {
        // Assumed XResources file
        theme = theme.fromXRFile(strings.Split(string(filedata), "\n"))
    } else {
        fmt.Println("File format (suffix) not valid.")
        os.Exit(1)
    }

    return theme
}

func set(theme XRTheme) {
    xdgConfig, err := os.UserConfigDir()
    check(err)

    tmpl, err := template.New("theme").Parse(
`{{- range $id, $value := .Entries -}}
{{$id}}: {{$value}}
{{end -}}`,
    )
    check(err)

    themefile, err := os.Create(xdgConfig + "/" + XRDB_FILENAME)
    check(err)
    defer themefile.Close()

    tmpl.Execute(themefile, theme)

    xrdbCmd := exec.Cmd{
        Dir: xdgConfig,
        Path: "/usr/bin/xrdb",
        Args: []string{ "-override", XRDB_FILENAME },
    }

    err = xrdbCmd.Run()
    check(err)
}

func save(theme XRTheme, filename string) {
    // Only pywal json for now
    var pwTheme PywalJsonTheme
    pwTheme = pwTheme.fromXRTheme(theme)
    jsonData, err := json.Marshal(pwTheme)
    check(err)

    err = ioutil.WriteFile(filename, jsonData, 0644)
    check(err)
}

func fileSuffix(filename string) string {
    slash := strings.LastIndex(filename, "/")
    if (slash >= 0) {
        filename = filename[slash:]
    }

    dot := strings.LastIndex(filename, ".")
    if dot >= 0 {
        return filename[dot + 1:]
    } else {
        return ""
    }
}

func (xrTheme XRTheme) fromXRFile(lines []string) XRTheme {
    xrTheme.Entries = make(map[string]string)

    i := 0
    for _, line := range lines {
        re_comment := regexp.MustCompile(`^\s*;`)
        if (re_comment.Match([]byte(line))) {
            continue
        }

        re_entry := regexp.MustCompile(
            `(?P<Id>^[\w\d\*]+\.[\w\d\*\.]+):\s*` + // e.g "*.color1: "
            `(?P<Value>[#\w\d\*\.]+\s*$)`,          // e.g "#EAF0B1"
        )

        submatches := re_entry.FindSubmatch([]byte(line))
        if (len(submatches) == 3) {
            // 3 values: [ full, id, value ]
            xrTheme.Entries[string(submatches[1])] = string(submatches[2])
            i++
        }
    }

    return xrTheme
}

func (pwTheme PywalJsonTheme) toXRTheme() XRTheme {
    theme := XRTheme {
        Entries: map[string]string {
            "*.foreground": pwTheme.Special.Foreground,
            "*.background": pwTheme.Special.Background,
            "*.cursor": pwTheme.Special.Cursor,
            "*.color0": pwTheme.Colors.Color0,
            "*.color1": pwTheme.Colors.Color1,
            "*.color2": pwTheme.Colors.Color2,
            "*.color3": pwTheme.Colors.Color3,
            "*.color4": pwTheme.Colors.Color4,
            "*.color5": pwTheme.Colors.Color5,
            "*.color6": pwTheme.Colors.Color6,
            "*.color7": pwTheme.Colors.Color7,
            "*.color8": pwTheme.Colors.Color8,
            "*.color9": pwTheme.Colors.Color9,
            "*.color10": pwTheme.Colors.Color10,
            "*.color11": pwTheme.Colors.Color11,
            "*.color12": pwTheme.Colors.Color12,
            "*.color13": pwTheme.Colors.Color13,
            "*.color14": pwTheme.Colors.Color14,
            "*.color15": pwTheme.Colors.Color15,
        },
    }

    return theme
}

func (pwTheme PywalJsonTheme) fromXRTheme(xrTheme XRTheme) PywalJsonTheme {
    pwTheme = PywalJsonTheme {
        Special: PywalJsonSpecial {
            Foreground: xrTheme.Entries["*.foreground"],
            Background: xrTheme.Entries["*.background"],
            Cursor: xrTheme.Entries["*.cursor"],
        },
        Colors: PywalJsonColors {
            Color0: xrTheme.Entries["*.color0"],
            Color1: xrTheme.Entries["*.color1"],
            Color2: xrTheme.Entries["*.color2"],
            Color3: xrTheme.Entries["*.color3"],
            Color4: xrTheme.Entries["*.color4"],
            Color5: xrTheme.Entries["*.color5"],
            Color6: xrTheme.Entries["*.color6"],
            Color7: xrTheme.Entries["*.color7"],
            Color8: xrTheme.Entries["*.color8"],
            Color9: xrTheme.Entries["*.color9"],
            Color10: xrTheme.Entries["*.color10"],
            Color11: xrTheme.Entries["*.color11"],
            Color12: xrTheme.Entries["*.color12"],
            Color13: xrTheme.Entries["*.color13"],
            Color14: xrTheme.Entries["*.color14"],
            Color15: xrTheme.Entries["*.color15"],
        },
    }

    return pwTheme
}
