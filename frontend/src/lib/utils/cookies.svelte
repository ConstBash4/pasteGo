<script lang="ts" context="module">
	import { refresh } from '$lib/api/auth/auth.svelte';
	let isRefreshing: boolean = false;

	export function getCookie(cookieName: string): string | null {
		const cookies = document.cookie.split('; ');

		for (const cookie of cookies) {
			const [key, value] = cookie.split('=');
			if (key === cookieName) {
				return decodeURIComponent(value);
			}
		}

		return null;
	}

	// Функция для проверки и обновления токенов
	export async function checkAndRefreshTokens() {
		if (isRefreshing) {
			return;
		}
		isRefreshing = true;
		console.log('Проверяю токен');
		// Проверяем наличие куки 'exp'
		const expCookie = getCookie('exp');
		if (!expCookie) {
			console.log('Exp не найден');
			isRefreshing = false;
			return;
		}

		// Парсим время истечения токена
		const expTime = parseInt(expCookie, 10);
		if (isNaN(expTime)) {
			return;
		}

		// Текущее время в секундах
		const now = Date.now() / 1000;
		const timeUntilExpire = expTime - now;

		// Если токен уже истёк
		/*if (timeUntilExpire <= 0) {
			console.log('Токен истек');
			isRefreshing = false;
			return;
		}*/

		// Если до истечения осталось менее 30 секунд
		if (timeUntilExpire <= 30 || timeUntilExpire <= 0) {
			try {
				console.log('Начал обновлять токен');
				// Выполняем запрос на обновление токенов
				await refresh();

				// После успешного обновления, проверяем снова
				isRefreshing = false;
				checkAndRefreshTokens();
			} catch (error) {
				console.error('Ошибка при обновлении токенов:', error);
			}
		} else {
			// Устанавливаем таймаут для следующей проверки
			console.log('Устанавливаю таймаут');
			isRefreshing = false;
			setTimeout(
				() => {
					checkAndRefreshTokens();
				},
				(timeUntilExpire - 30) * 1000
			);
		}
	}
</script>
