<script>
  import { DialogContent } from "svelte-dialogs";
  import * as wailsRuntime from "../wailsjs/runtime/runtime";

  export let version;
</script>

<DialogContent>
  <h1 slot="header" style="color:black;">{version.tag_name}</h1>
  <svelte:fragment slot="body">
    <p style="color:black;min-height:50%">{@html version.body}</p>
  </svelte:fragment>
  <svelte:fragment slot="footer">
    <button
      style="height:34px;width:80%;"
      on:click={() => {
        wailsRuntime.BrowserOpenURL(version.html_url);
      }}>立即更新</button
    >
  </svelte:fragment>
</DialogContent>
