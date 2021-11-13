<script lang="ts">
	export let src: string;
	export let alt: string;
</script>

<figure>
	<img {src} {alt} height={48} width={48} />
</figure>
