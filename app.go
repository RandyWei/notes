package main

import (
	"context"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	content  string //markdown 内容
	fileName string //文件名
	filePath string //文件路径
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetFileName() string {
	return a.fileName
}

func (a *App) ResizeWindows() {
	runtime.WindowToggleMaximise(a.ctx)
}

// 当 Vditor 编辑改变时，返回数据，目前只有这种方式才能从 frontend 拿到数据
func (a *App) OnVditorChanged(value string) {
	a.content = value
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
}
