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
	content  string
	fileName string
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

func (a *App) OnVditorChanged(value string) {
	a.content = value
}

// 保存文件，
// filePath为文件路径
func (a *App) saveFile(filePath string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	defer file.Close()

	_, err = file.WriteString(a.content)

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
