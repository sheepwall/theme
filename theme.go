package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
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

// XResources line
type XREntry struct {
    Id string
    Value string
}

// XReasources theme (several lines)
type XRTheme struct {
    Entries []XREntry
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    xdgConfig, err := os.UserConfigDir()
    check(err)

    if (len(os.Args) < 2) {
        fmt.Println("Invalid usage: Supply a filename.")
        os.Exit(1)
    }
    filename := os.Args[1]
    fmt.Println("Reading file ", filename)
    filedata, err := ioutil.ReadFile(filename)
    check(err)

    var pwTheme PywalJsonTheme
    err = json.Unmarshal(filedata, &pwTheme)
    check(err)

    theme := pwTheme.toXRTheme()

    tmpl, err := template.New("theme").Parse(
`{{- range .Entries -}}
{{.Id}}: {{.Value}}
{{end -}}`,
    )
    check(err)

    themefile, err := os.Create(xdgConfig + "/Xtheme")
    check(err)
    defer themefile.Close()

    tmpl.Execute(themefile, theme)

    xrdbCmd := exec.Cmd{
        Dir: xdgConfig,
        Path: "/usr/bin/xrdb",
        Args: []string{"-override", "Xtheme"},
    }

    err = xrdbCmd.Run()
    check(err)
}

func (pwTheme PywalJsonTheme) toXRTheme() XRTheme {
    theme := XRTheme {
        Entries: []XREntry {
            {Id: "*.foreground", Value: pwTheme.Special.Foreground },
            {Id: "*.background", Value: pwTheme.Special.Background },
            {Id: "*.cursor", Value: pwTheme.Special.Cursor },
            {Id: "*.color0", Value: pwTheme.Colors.Color0 },
            {Id: "*.color1", Value: pwTheme.Colors.Color1 },
            {Id: "*.color2", Value: pwTheme.Colors.Color2 },
            {Id: "*.color3", Value: pwTheme.Colors.Color3 },
            {Id: "*.color4", Value: pwTheme.Colors.Color4 },
            {Id: "*.color5", Value: pwTheme.Colors.Color5 },
            {Id: "*.color6", Value: pwTheme.Colors.Color6 },
            {Id: "*.color7", Value: pwTheme.Colors.Color7 },
            {Id: "*.color8", Value: pwTheme.Colors.Color8 },
            {Id: "*.color9", Value: pwTheme.Colors.Color9 },
            {Id: "*.color10", Value: pwTheme.Colors.Color10 },
            {Id: "*.color11", Value: pwTheme.Colors.Color11 },
            {Id: "*.color12", Value: pwTheme.Colors.Color12 },
            {Id: "*.color13", Value: pwTheme.Colors.Color13 },
            {Id: "*.color14", Value: pwTheme.Colors.Color14 },
            {Id: "*.color15", Value: pwTheme.Colors.Color15 },
        },
    }

    return theme
}
