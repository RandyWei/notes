package main

import (
	"context"
	"fmt"
	"os"
	goRuntime "runtime"
	"strings"

	"icu.bughub.app/notes/backend/network"
	"icu.bughub.app/notes/backend/repo"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	content  string //markdown 内容
	fileName string //文件名
	filePath string //文件路径
	saved    bool   //当前是否已经保存
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.saved = true
	a.ctx = ctx

}

func (a *App) onDomReady(ctx context.Context) {
	a.update()
	runtime.WindowSetTitle(a.ctx, a.fileName)
}

func (a *App) GetFileName() string {
	return a.fileName
}

func (a *App) ResizeWindow() {
	runtime.WindowToggleMaximise(a.ctx)
}

// 当 Vditor 编辑改变时，返回数据，目前只有这种方式才能从 frontend 拿到数据
func (a *App) OnVditorValueChanged(value string) {
	a.content = value
	a.saved = false
}

func (a *App) LoadContentFromLocal() string {
	return a.content
}

func (a *App) OS() string {
	return goRuntime.GOOS
}

// 保存文件
func (a *App) saveFile() {
	file, err := os.OpenFile(a.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	defer file.Close()

	_, err = file.WriteString(a.content)

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	a.saved = true
}

var CurrentVersion = "0.0.1"

func (a *App) update() {
	//https://api.github.com/repos/RandyWei/notes/releases/latest
	//https://api.github.com/repos/alist-org/alist/releases/latest
	request := network.Request{
		Url: "https://api.github.com/repos/RandyWei/notes/releases/latest",
	}
	var lastest repo.Release
	message, err := request.PostParse(&lastest)
	if err != nil {
		fmt.Printf("message: %v\n", message)
		fmt.Printf("err: %v\n", err)
		return
	}

	version := strings.Split(strings.Replace(lastest.TagName, "v", "", 1), ".")

	if len(version) < 3 {
		return
	}
	major := version[0]
	minor := version[1]
	revision := version[2]

	currentVersion := strings.Split(CurrentVersion, ".")

	//判断版本号
	if revision <= currentVersion[2] {
		if minor <= currentVersion[1] {
			if major <= currentVersion[0] {
				return
			}
		}
	}

	//判断 assets 中是否包含对应平台的文件
	for _, asset := range lastest.Assets {
		//application/x-diskcopy
		if (asset.ContentType == "application/gzip" && goRuntime.GOOS == "darwin") || (asset.ContentType == "application/x-msdownload" && goRuntime.GOOS == "windows") { //dmg文件
			runtime.EventsEmit(a.ctx, "OnUpdate", lastest)
			break
		}
	}

}
