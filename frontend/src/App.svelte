<script lang="ts">
  import { onMount, afterUpdate } from "svelte";
  import * as wailsRuntime from "../wailsjs/runtime/runtime";

  import Vditor from "vditor";
  import {
    OnVditorChanged,
    ResizeWindows,
    GetFileName,
  } from "../wailsjs/go/main/App";

  let titleBarHeight = 35;

  let editor;

  let fileName = "";

  function resizeWindows() {
    ResizeWindows();
  }

  onMount(() => {
    //初始化 vditor
    editor = new Vditor("vditor", {
      width: "100%",
      height: "calc(100vh - 35px)", //最好设置父布局 main 的 height 为100vh
      theme: "dark",
      icon: "material",
      preview: {
        theme: { current: "dark" },
      },
      toolbarConfig: {
        pin: true,
      },
    });
    //获取 fileName
    GetFileName().then((name) => {
      fileName = name;
    });
    //监听保存时文件名的变化
    wailsRuntime.EventsOn("OnFileNameChanged", (name) => {
      console.log(name);
      fileName = name;
    });
  });
</script>

<main>
  <div
    id="title"
    style="--wails-draggable:drag;;height:{titleBarHeight}px;"
    on:dblclick={resizeWindows}
  >
    {fileName}
  </div>
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
