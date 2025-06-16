<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import { onMount } from 'svelte';
	import { getCookie } from '$lib/utils/cookies.svelte';
	import { goto } from '$app/navigation';
	import { z } from 'zod';
	import Frame from '$lib/components/Frame.svelte';
	import {
		PasteInfoListSchema,
		PasteInfoSchema,
		type PasteInfo,
		type PasteInfoList
	} from '$lib/api/types.svelte';
	import { createPaste, deletePaste, getPasteList, updatePaste } from '$lib/api/paste/paste.svelte';
	import { formatUnixTime } from '$lib/utils/time.svelte';

	let mode: 'auth' | 'profile' | '' = 'auth';
	let username: string = '';
	let error: string | null = null;
	let showForm = false;
	let pasteText = '';
	let isPublic = false;
	let hasPassword = false;
	let password = '';
	let timeToLive = 'forever';
	let pasteList: PasteInfoList = {
		pastes: []
	};
	let currentPaste: PasteInfo | null = null;
	let editingPasteId: string | null = null;

	const timeOptions = [
		{ value: 'minute', label: 'Минута' },
		{ value: 'hour', label: 'Час' },
		{ value: 'day', label: 'День' },
		{ value: 'week', label: 'Неделя' },
		{ value: 'month', label: 'Месяц' },
		{ value: 'year', label: 'Год' },
		{ value: 'forever', label: 'Навсегда' }
	];

	onMount(async () => {
		const usernameCookie = getCookie('username');
		if (usernameCookie) {
			mode = 'profile';
			username = usernameCookie;
		} else {
			goto('/auth/login');
		}

		try {
			let response = await getPasteList();
			if (response.code == 0) {
				pasteList = PasteInfoListSchema.parse(response.message);
			} else {
				error =
					response.code + ': ' + response.explanation ||
					'Произошла ошибка при получении списка вставок';
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
	});

	async function handleCreatePaste() {
		error = null;
		try {
			// Валидация данных перед отправкой
			const validatedData = PasteInfoSchema.parse({
				id: '',
				author: username,
				created: 0,
				updated: 0,
				expTime: 0,
				lifetime: timeToLive,
				text: pasteText,
				password: password,
				hasPassword: hasPassword,
				public: isPublic
			});

			let response = await createPaste(validatedData);

			if (response.code == 0) {
				let createdPaste: PasteInfo = PasteInfoSchema.parse(response.message);
				//pasteList.pastes.unshift(createdPaste);
				pasteList.pastes = [createdPaste, ...pasteList.pastes];
				showForm = false;
			} else {
				error =
					response.code + ': ' + response.explanation || 'Произошла ошибка при создании вставки';
			}
		} catch (err) {
			// Обработка ошибок валидации или других ошибок
			if (err instanceof z.ZodError) {
				error = err.issues[0]?.message || 'Ошибка валидации';
			} else {
				error = 'Произошла ошибка при создании вставки:' + err;
			}
		}
	}

	async function handleDeletePaste(id: string) {
		error = null;
		try {
			let response = await deletePaste(id);

			if (response.code == 0) {
				deleteFromPasteList(id);
			} else {
				error =
					response.code + ': ' + response.explanation || 'Произошла ошибка при удалении вставки';
			}
		} catch (err) {
			// Обработка ошибок валидации или других ошибок
			if (err instanceof z.ZodError) {
				error = err.issues[0]?.message || 'Ошибка валидации';
			} else {
				error = 'Произошла ошибка при удалении вставки:' + err;
			}
		}
	}

	function deleteFromPasteList(id: string) {
		pasteList.pastes = pasteList.pastes.filter((paste) => paste.id !== id);
	}

	function handleEditPaste(paste: PasteInfo) {
		currentPaste = { ...paste };
		editingPasteId = paste.id;
	}

	function handleCancelEdit() {
		currentPaste = null;
		editingPasteId = null;
	}

	async function handleSavePaste() {
		console.log('UpdatedPaste');
		console.log(currentPaste);
		if (!currentPaste) return;
		error = null;
		try {
			// Валидация данных перед отправкой
			//const validatedData = PasteInfoSchema.parse({ currentPaste });

			let response = await updatePaste(currentPaste.id, currentPaste);

			if (response.code == 0) {
				let updatedPaste: PasteInfo = PasteInfoSchema.parse(response.message);
				// Обновление массива заметок
				pasteList.pastes = pasteList.pastes.map((paste) =>
					paste.id === updatedPaste.id ? updatedPaste : paste
				);

				// Сброс состояния
				currentPaste = null;
				editingPasteId = null;
			} else {
				error =
					response.code + ': ' + response.explanation || 'Произошла ошибка при изменении вставки';
			}
		} catch (err) {
			// Обработка ошибок валидации или других ошибок
			if (err instanceof z.ZodError) {
				error = err.issues[0]?.message || 'Ошибка валидации';
			} else {
				error = 'Произошла ошибка при изменении вставки:' + err;
			}
		}
	}
</script>

<svelte:head>
	<title>Pastes</title>
</svelte:head>

<Header {mode} {username} />
<Frame {content} />
<Footer />

{#snippet content()}
	<button class="showForm" on:click={() => (showForm = !showForm)}> Создать вставку </button>
	{#if error}
		<div class="error">{error}</div>
	{/if}
	{#if showForm}
		<form on:submit|preventDefault={handleCreatePaste} class="container">
			<div class="form-group">
				<label for="pasteText">Текст вставки</label>
				<textarea
					id="pasteText"
					bind:value={pasteText}
					required
					placeholder="Введите текст вставки"
					class="textarea-control"
				></textarea>
			</div>

			<div class="checkbox-group">
				<label>
					Публичная
					<input type="checkbox" class="checkbox-input" bind:checked={isPublic} />
				</label>

				<label>
					Пароль
					<input type="checkbox" class="checkbox-input" bind:checked={hasPassword} />
				</label>

				{#if hasPassword}
					<div class="form-group">
						<input
							type="password"
							id="password"
							bind:value={password}
							required
							placeholder="Введите пароль"
						/>
					</div>
				{/if}
			</div>

			<div class="form-group">
				<label for="timeToLive">Время жизни</label>
				<select bind:value={timeToLive} id="timeToLive">
					{#each timeOptions as option}
						<option value={option.value}>{option.label}</option>
					{/each}
				</select>
			</div>

			<button type="submit">Создать заметку</button>
		</form>
	{/if}
	{#if pasteList.pastes.length !== 0}
		{#each pasteList.pastes as paste}
			{#if paste.id === editingPasteId && currentPaste !== null}
				<form on:submit|preventDefault={handleSavePaste} class="container">
					<button on:click={() => handleCancelEdit()}>Отмена</button>
					<div class="form-group">
						<label for="pasteText">Текст вставки</label>
						<textarea
							id="pasteText"
							bind:value={currentPaste.text}
							required
							placeholder="Введите текст вставки"
							class="textarea-control"
						></textarea>
					</div>

					<div class="checkbox-group">
						<label>
							Публичная
							<input type="checkbox" class="checkbox-input" bind:checked={currentPaste.public} />
						</label>
						<label>
							Пароль
							<input
								type="checkbox"
								class="checkbox-input"
								bind:checked={currentPaste.hasPassword}
							/>
						</label>
						{#if currentPaste.hasPassword}
							<div class="form-group">
								<input
									type="password"
									id="password"
									bind:value={currentPaste.password}
									placeholder="Введите новый пароль (если требуется)"
								/>
							</div>
						{/if}
					</div>

					<div class="form-group">
						<label for="timeToLive">Время жизни</label>
						<select bind:value={currentPaste.lifetime} id="timeToLive">
							{#each timeOptions as option}
								<option value={option.value}>{option.label}</option>
							{/each}
						</select>
					</div>

					<button type="submit">Обновить</button>
				</form>
			{:else}
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
						<div class="paste-date">
							<strong>Дата истечения:</strong> <span>{formatUnixTime(paste.expTime)}</span>
						</div>
						<div class="paste-meta">
							{#if paste.hasPassword}
								<span>🔒 Наличие пароля: <strong>Да</strong></span><br />
							{:else}
								<span>🔒 Наличие пароля: <strong>Нет</strong></span><br />
							{/if}
							{#if paste.public}
								<span>🌐 Публичность: <strong>Да</strong></span>
							{:else}
								<span>🌐 Публичность: <strong>Нет</strong></span>
							{/if}
						</div>
					</div>

					<!-- Текст заметки -->
					<div class="paste-text">
						<p>
							{paste.text}
						</p>
						<p>http://localhost:10015/paste/{paste.id}</p>
					</div>
					<div class="paste-actions">
						<button on:click={() => handleEditPaste(paste)}>Редактировать</button>
						<button on:click={() => handleDeletePaste(paste.id)}>Удалить</button>
					</div>
				</div>
			{/if}
		{/each}
	{/if}
{/snippet}

<style>
	.container {
		background-color: #1e1e1e;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
		max-width: 400px;
		width: 100%;
		margin: 10px auto;
	}

	.checkbox-input {
		accent-color: #00ffcc;
	}

	.error {
		color: #ff4757;
		font-size: 1rem;
		margin-top: 1rem;
		text-align: center;
	}

	.paste-actions {
		display: flex;
		gap: 12px;
	}

	.paste-info {
		gap: 10px;
		display: flex;
		flex-direction: column;
	}

	.container textarea,
	input,
	select {
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

	.showForm {
		margin: 10px auto;
		max-width: max-content;
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-group textarea:focus,
	input:focus,
	select:focus {
		border-color: #00ffcc;
		outline: none;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		color: #ccc;
	}

	.form-group textarea,
	input,
	select {
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
</style>
