<script lang="ts">
  import { onMount, afterUpdate } from "svelte";
  import * as wailsRuntime from "../wailsjs/runtime/runtime";
  import ContextMenu from "./ContextMenu.svelte";
  import { dialogs } from "svelte-dialogs";
  import Vditor from "vditor";
  import {
    OnVditorValueChanged,
    ResizeWindow,
    GetFileName,
  } from "../wailsjs/go/main/App";

  let titleBarHeight = 35;

  let editor;

  let fileName = "";

  let options: IOptions = {
    width: "100%",
    height: "calc(100vh - 35px)", //最好设置父布局 main 的 height 为100vh
    theme: "dark",
    icon: "material",
    preview: {
      theme: { current: "dark" },
    },
    cache: {
      enable: false,
    },
    toolbar: [
      "headings",
      "bold",
      "strike",
      "link",
      "list",
      "ordered-list",
      "check",
      "outdent",
      "indent",
      "quote",
      "line",
      "table",
      "export",
      "outline",
    ],
    toolbarConfig: {
      pin: true,
    },
    input(value) {
      OnVditorValueChanged(value);
    },
  };

  function resizeWindow() {
    ResizeWindow();
  }

  onMount(() => {
    //初始化 vditor
    editor = new Vditor("vditor", options);
    //获取 fileName
    GetFileName().then((name) => {
      fileName = name;
    });
    //监听保存时文件名的变化
    wailsRuntime.EventsOn("OnFileNameChanged", (name) => {
      fileName = name;
    });
    //监听菜单点击
    wailsRuntime.EventsOn("OnMenuClick", (menuName) => {
      if (menuName == "about") {
        const htmlString = `
        <div>
            <h1 id="dialog-title-id">这里什么都没有</h1>
        </div>`;
        const opts = {
          closeButton: false,
        };
        dialogs.modal(htmlString, opts);
      }
    });
    //监听文件内容加载
    wailsRuntime.EventsOn("OnLoadFile", (content) => {
      editor.setValue(content);
    });
  });
</script>

<main>
  <div
    id="title"
    style="--wails-draggable:drag;;height:{titleBarHeight}px;"
    on:dblclick={resizeWindow}
  >
    {fileName}
  </div>
  <!-- <ContextMenu /> -->
  <div id="vditor" />
</main>

<style>
  #title {
    display: flex;
    align-items: center;
    justify-content: center;
    user-select: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    background-color: #1d2125;
  }
</style>
