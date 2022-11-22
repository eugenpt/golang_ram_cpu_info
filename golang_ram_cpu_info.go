package main

import (
    "fmt"
    "io/ioutil"
    "time"

    "github.com/getlantern/systray"
    "github.com/getlantern/systray/example/icon"
)

//     "github.com/skratchdot/open-golang/open"

func main() {
    fmt.Println("Hello, World!")

    onExit := func() {
        now := time.Now()
        ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
    }

    systray.Run(onReady, onExit)
}

func onReady() {
    systray.SetTemplateIcon(icon.Data, icon.Data)
    systray.SetTitle("Awesome App")
    systray.SetTooltip("Lantern")
    mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
    go func() {
        <-mQuitOrig.ClickedCh
        fmt.Println("Requesting quit")
        systray.Quit()
        fmt.Println("Finished quitting")
    }()
}
