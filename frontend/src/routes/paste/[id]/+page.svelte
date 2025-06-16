<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import Frame from '$lib/components/Frame.svelte';
	import { getCookie } from '$lib/utils/cookies.svelte';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { getPaste } from '$lib/api/paste/paste.svelte';
	import { z } from 'zod';
	import {
		PasteInfoSchema,
		PastePasswordSchema,
		type PasteInfo,
		type PastePassword
	} from '$lib/api/types.svelte';
	import { formatUnixTime } from '$lib/utils/time.svelte';
	let mode: 'auth' | 'profile' | '' = 'auth';
	let username: string = '';
	let error: string | null = null;
	let passwordRequired = false;
	let psw: PastePassword = PastePasswordSchema.parse({ password: '' });
	let paste: PasteInfo | null = null;
	onMount(async () => {
		const usernameCookie = getCookie('username');
		if (usernameCookie) {
			mode = 'profile';
			username = usernameCookie;
		}
		await handleGetPaste();
	});

	async function handleGetPaste() {
		error = null;
		try {
			let response = await getPaste(page.params.id, psw);
			switch (response.code) {
				case 0: {
					paste = PasteInfoSchema.parse(response.message);
					passwordRequired = false;
					break;
				}
				case 2004: {
					passwordRequired = true;
					break;
				}
				default: {
					error =
						response.code + ': ' + response.explanation ||
						'Произошла ошибка при получении списка вставок';
				}
			}
		} catch (err) {
			// Обработка ошибок валидации или других ошибок
			if (err instanceof z.ZodError) {
				error = err.issues[0]?.message || 'Ошибка валидации';
				console.log('Ощибка Zod:' + error);
			} else {
				error = 'Произошла ошибка при получении списка вставок:' + err;
			}
		}
	}
</script>

<svelte:head>
	<title>Paste</title>
</svelte:head>

<Header {mode} {username} />
<Frame {content} />
<Footer />

{#snippet content()}
	{#if error}
		<div class="error">{error}</div>
	{/if}
	{#if passwordRequired}
		<form on:submit|preventDefault={handleGetPaste} class="container">
			<div class="form-group">
				<label>
					Эта вставка требует пароль
					<input
						type="password"
						id="password"
						bind:value={psw.password}
						required
						placeholder="Введите пароль"
					/>
				</label>
				<button type="submit">Подтвердить</button>
			</div>
		</form>
	{/if}
	{#if paste != null}
		<div class="container">
			<!-- Информация о заметке -->
			<div class="paste-info">
				<div class="paste-date">
					<strong>Дата создания:</strong> <span>{formatUnixTime(paste.created)}</span>
				</div>
				{#if paste.updated !== -1 && paste.updated != undefined}
					<div class="paste-date">
						<strong>Дата обновления:</strong> <span>{formatUnixTime(paste.updated)}</span>
					</div>
				{/if}
				<div class="paste-meta">
					<span>🗿 <strong>Автор:</strong> {paste.author}</span>
				</div>
			</div>

			<!-- Текст заметки -->
			<div class="paste-text">
				<p>
					{paste.text}
				</p>
			</div>
		</div>
	{/if}
{/snippet}

<style>
	.error {
		color: #ff4757;
		font-size: 1rem;
		margin-top: 1rem;
		text-align: center;
	}

	.container {
		background-color: #1e1e1e;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
		max-width: 400px;
		width: 100%;
		margin: 10px auto;
	}

	.paste-info {
		gap: 10px;
		display: flex;
		flex-direction: column;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-group input:focus {
		border-color: #00ffcc;
		outline: none;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		color: #ccc;
	}

	.form-group input {
		width: 100%;
		padding: 0.75rem;
		background-color: #2a2a2a;
		border: 1px solid #444;
		color: #fff;
		border-radius: 5px;
		font-size: 1rem;
		transition: border-color 0.3s;
		box-sizing: border-box;
		resize: vertical;
		min-height: 45px;
	}

	.container input {
		width: 100%;
		padding: 0.5rem;
		font-size: 1rem;
		margin-top: 0.5rem;
	}

	button {
		width: 100%;
		padding: 0.75rem;
		background-color: #00ffcc;
		margin-top: 10px;
		color: #000;
		border: none;
		border-radius: 5px;
		font-size: 1rem;
		cursor: pointer;
		transition: background-color 0.3s;
	}

	button:hover {
		background-color: #00e0b8;
	}

	button:disabled {
		background-color: #00e0b8;
		cursor: not-allowed;
	}
</style>
