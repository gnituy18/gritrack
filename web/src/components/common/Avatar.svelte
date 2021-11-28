<script lang="ts">
	export let src: string;
	export let alt: string;
</script>

<figure>
	<img class="w-10 rounded-full" {src} {alt} />
</figure>
