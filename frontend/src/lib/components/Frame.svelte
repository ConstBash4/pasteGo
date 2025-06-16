<script lang="ts">
	let { content } = $props();
</script>

<div class="frame">
	{@render content?.()}
</div>

<style>
	.frame {
		width: auto;
		height: 100%;
		background-color: #121212;
		display: flex;
		flex-direction: column;
		overflow-y: auto;
		margin-top: 64px;
		flex: 1;
	}
</style>
