package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	app.fileName = "未命名.md"

	//总菜单
	notesMenu := menu.NewMenu()
	//appMenu，只针对 Mac
	if runtime.GOOS == "darwin" {
		notesMenu.Append(menu.AppMenu())
	}

	fileSubMenu := notesMenu.AddSubmenu("文件")
	fileSubMenu.AddText("新建文件", keys.CmdOrCtrl("n"), func(cd *menu.CallbackData) {
		//打开文件之前，判断一下当前打开的是否已经保存，目前不支持多窗口
		if !app.saved {
			wailsRuntime.MessageDialog(app.ctx, wailsRuntime.MessageDialogOptions{
				Type:    wailsRuntime.WarningDialog,
				Title:   "警告",
				Message: "当前内容尚未保存，请保存后再打开新文件",
			})
			return
		}
		app.filePath = ""
		app.fileName = "未命名.md"
		app.content = ""

		wailsRuntime.EventsEmit(app.ctx, "OnFileNameChanged", app.fileName)
		wailsRuntime.EventsEmit(app.ctx, "OnLoadFile", app.content)
	})
	fileSubMenu.AddText("打开文件", keys.CmdOrCtrl("o"), func(cd *menu.CallbackData) {

		//打开文件之前，判断一下当前打开的是否已经保存，目前不支持多窗口
		if !app.saved {
			wailsRuntime.MessageDialog(app.ctx, wailsRuntime.MessageDialogOptions{
				Type:    wailsRuntime.WarningDialog,
				Title:   "警告",
				Message: "当前内容尚未保存，请保存后再打开新文件",
			})
			return
		}

		filePath, err := wailsRuntime.OpenFileDialog(app.ctx, wailsRuntime.OpenDialogOptions{Title: "打开文件"})
		if err != nil {
			fmt.Printf("err: %T\n", err)
		}
		app.filePath = filePath
		//分割出目录和文件名
		_, fileName := filepath.Split(filePath)
		if fileName != "" {
			app.fileName = fileName
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("err: %T\n", err)
		}

		app.content = string(data)

		wailsRuntime.EventsEmit(app.ctx, "OnFileNameChanged", app.fileName)
		wailsRuntime.EventsEmit(app.ctx, "OnLoadFile", app.content)
	})
	fileSubMenu.AddText("保存文件", keys.CmdOrCtrl("s"), func(cd *menu.CallbackData) {

		//如果文件不存在，则弹窗选择文件
		//有两种情况：1、新建文件，从未保存过；2、从已有文件打开后，再删除本地文件
		if _, err := os.Stat(app.filePath); err != nil && os.IsNotExist(err) {
			filePath, err := wailsRuntime.SaveFileDialog(app.ctx, wailsRuntime.SaveDialogOptions{
				Title:                "保存文件",
				DefaultFilename:      app.fileName,
				CanCreateDirectories: true,
			})

			if err != nil {
				fmt.Printf("err: %T\n", err)
			}
			app.filePath = filePath
			//分割出目录和文件名
			_, fileName := filepath.Split(filePath)
			if fileName != "" {
				app.fileName = fileName
			}
		}
		//如果用户手动删除掉后缀名，需要补全
		if !strings.HasSuffix(app.filePath, ".md") {
			app.filePath = fmt.Sprintf("%v.md", app.filePath)
		}
		wailsRuntime.EventsEmit(app.ctx, "OnFileNameChanged", app.fileName)
		app.saveFile()

	})
	//导出功能
	// fileSubMenu.AddSeparator()
	// exportMenu := fileSubMenu.AddSubmenu("导出")
	// exportMenu.AddText("PDF", nil, func(cd *menu.CallbackData) {})
	// exportMenu.AddText("HTML", nil, func(cd *menu.CallbackData) {})

	if runtime.GOOS == "darwin" {
		notesMenu.Append(menu.EditMenu())
	}

	//帮助
	helpSubMenu := notesMenu.AddSubmenu("帮助")
	helpSubMenu.AddText("关于", nil, func(cd *menu.CallbackData) {
		wailsRuntime.EventsEmit(app.ctx, "OnMenuClick", "about")
	})

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "码札",
		Width:     1024,
		Height:    768,
		MinWidth:  1024,
		MinHeight: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnDomReady:       app.onDomReady,
		Menu:             notesMenu,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHidden(),
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
