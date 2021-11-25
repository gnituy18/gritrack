<script lang="ts">
	export let src: string;
	export let alt: string;
</script>

<figure>
	<img {src} {alt} height={36} width={36} />
</figure>

<style>
	figure {
		margin: 0;
		height: 36px;
		width: 36px;
	}

	img {
		border-radius: 50%;
	}
</style>
