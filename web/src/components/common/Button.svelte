<script lang="ts">
  export let href: string = "";
  export let value: string = "";
  export let size: "s" | "m" | "l" = "s";
  export let theme: "primary" | "hidden" = "primary";
  export let onClick: () => Promise<void> | void = () => {};

  const pClass = (size: "s" | "m" | "l") => {
    switch (size) {
      case "s":
        return "px-3 py-1";
      case "m":
        return "px-5 py-2";
    }
  };

  const themeClass = (theme: "primary" | "hidden") => {
    switch (theme) {
      case "primary":
        return "btn";
      case "hidden":
        return "btn-hidden";
    }
  };
</script>

<a class="flex-none" {href}>
  <button class={`${themeClass(theme)} ${pClass(size)}`} on:click={onClick}>
    {value}
  </button>
</a>
