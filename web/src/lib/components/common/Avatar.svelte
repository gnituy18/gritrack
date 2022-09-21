<script lang="ts">
  export let src: string;
  export let alt: string;

  function changePictureRes(src: string, pixel: number): string {
    return src.replace(/s96-c/, `s${pixel}-c`);
  }
</script>

<figure>
  <img
    class="mx-auto w-24 rounded-full"
    src={changePictureRes(src, 256)}
    {alt}
  />
</figure>
