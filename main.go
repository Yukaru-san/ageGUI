package main

import (
	"flag"
	"log"
	"os"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

// Vars injected via ldflags by bundler
var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

// Vars used by the program itself
var (
	Window *astilectron.Window
	Logger *log.Logger
	fs     = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	Debug  = fs.Bool("d", false, "enables the debug mode") // CTRL-D -> Dev Tools
)

func main() {
	// Create Logger
	Logger = log.New(log.Writer(), log.Prefix(), log.Flags())
	/* Logging to a file instead
	f, _ := os.OpenFile("testlogfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	Logger.SetOutput(f)
	*/

	// Parse flags
	fs.Parse(os.Args[1:])

	//Run bootstrap
	err := createWindow()
	if err != nil {
		Logger.Fatal(err.Error())
	}

}

func createWindow() error {
	err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:         *Debug,
		Logger:        Logger,
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: HandleMessage,
			Adapter: func(i *astilectron.Window) {
				Window = i
			},
			Options: &astilectron.WindowOptions{
				Center:          astikit.BoolPtr(true),
				MinHeight:       astikit.IntPtr(710),
				Height:          astikit.IntPtr(710),
				MinWidth:        astikit.IntPtr(700),
				Width:           astikit.IntPtr(700),
				AutoHideMenuBar: astikit.BoolPtr(true),
			},
		}},
	})

	return err
}
