<script>
	import { checkAndRefreshTokens } from '$lib/utils/cookies.svelte';
	import { onMount } from 'svelte';

	let { children } = $props();

	// Запуск проверки при монтировании компонента
	onMount(() => {
		checkAndRefreshTokens();
	});
</script>

{@render children()}

<style>
	:global {
		body {
			margin: 0;
			padding: 0;
			font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
			background-color: #121212;
			color: #ffffff;
			display: flex;
			flex-direction: column;
			width: 100vw;
			height: 100vh;
		}
	}
</style>
