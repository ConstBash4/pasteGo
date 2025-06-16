<script lang="ts" context="module">
	import { apiClient } from '../api.svelte';
	import {
		APIResponseSchema,
		PasteInfoSchema,
		PastePasswordSchema,
		type APIResponse,
		type PasteInfo,
		type PastePassword
	} from '../types.svelte';

	export async function getPaste(id: string, data: PastePassword): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/paste/' + id,
			method: 'POST',
			requestData: data,
			requestSchema: PastePasswordSchema,
			responseSchema: APIResponseSchema
		});
	}

	export async function getPasteList(): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/v1/paste',
			method: 'GET',
			responseSchema: APIResponseSchema
		});
	}

	export async function createPaste(data: PasteInfo): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/v1/paste',
			method: 'POST',
			requestData: data,
			requestSchema: PasteInfoSchema,
			responseSchema: APIResponseSchema
		});
	}

	export async function updatePaste(id: string, data: PasteInfo): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/v1/paste/' + id,
			method: 'PUT',
			requestData: data,
			requestSchema: PasteInfoSchema,
			responseSchema: APIResponseSchema
		});
	}

	export async function deletePaste(id: string): Promise<APIResponse> {
		return apiClient.fetch({
			url: '/v1/paste/' + id,
			method: 'DELETE',
			responseSchema: APIResponseSchema
		});
	}
</script>
