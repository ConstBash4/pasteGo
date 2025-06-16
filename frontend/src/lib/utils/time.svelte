<script lang="ts" context="module">
	export function formatUnixTime(unixSec: number): string {
		if (unixSec === -1) {
			return 'вечность';
		}
		const date = new Date(unixSec * 1000); // Переводим секунды в миллисекунды

		// Форматируем время в HH:MM:SS
		const timeFormat = date.toLocaleTimeString('en-US', {
			hour12: false,
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		});

		// Форматируем дату в DD-MM-YYYY
		const dateFormat = date.toLocaleDateString('en-US', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric'
		});

		return `${timeFormat} ${dateFormat}`;
	}
</script>
