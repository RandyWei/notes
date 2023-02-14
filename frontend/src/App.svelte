<script lang="ts">
  import { onMount, afterUpdate } from "svelte";
  import * as wailsRuntime from "../wailsjs/runtime/runtime";

  import Vditor from "vditor";
  import { OnVditorChanged, ResizeWindows } from "../wailsjs/go/main/App";

  let titleBarHeight = 35;

  let innerHeight;

  let editor;

  //监听 innerHeight 变化，则刷新编辑器
  $: editor = new Vditor("vditor", {
    width: "100%",
    height: innerHeight - titleBarHeight,
    theme: "dark",
    icon: "material",
    preview: {
      theme: { current: "dark" },
    },
    toolbarConfig: {
      pin: true,
    },
    input(value) {
      OnVditorChanged(value);
    },
  });

  (<any>window).getVditorValue = () => {
    return editor.getValue();
  };

  function resizeWindows() {
    // wailsRuntime.WindowToggleMaximise();
    ResizeWindows();
  }

  onMount(() => {
    editor = new Vditor("vditor", {
      width: "100%",
      height: innerHeight - titleBarHeight,
      theme: "dark",
      icon: "material",
      preview: {
        theme: { current: "dark" },
      },
      toolbarConfig: {
        pin: true,
      },
    });
  });
</script>

<svelte:window bind:innerHeight />
<main>
  <div
    id="title"
    style="--wails-draggable:drag;;height:{titleBarHeight}px;"
    on:dblclick={resizeWindows}
  >
    标题
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
