<script lang="ts" context="module">
	import axios, {
		AxiosError,
		type AxiosInstance,
		type AxiosRequestConfig,
		type AxiosResponse
	} from 'axios';
	import { z } from 'zod';

	type FetchOptions<TRequest, TResponse> = {
		url: string;
		method: 'GET' | 'POST' | 'PUT' | 'DELETE';
		requestData?: TRequest;
		requestSchema?: z.ZodSchema<TRequest>;
		responseSchema: z.ZodSchema<TResponse>;
		config?: AxiosRequestConfig;
	};

	class ApiClient {
		private axiosInstance: AxiosInstance;

		constructor(baseURL: string) {
			this.axiosInstance = axios.create({
				baseURL,
				timeout: 10000,
				withCredentials: true
			});
		}

		// Метод для запросов
		async fetch<TRequest, TResponse>(
			options: FetchOptions<TRequest, TResponse>
		): Promise<TResponse> {
			try {
				// Валидация входных данных
				if (options.requestData && options.requestSchema) {
					options.requestData = options.requestSchema.parse(options.requestData);
				}

				// Выполнение запроса
				const response: AxiosResponse<TResponse> = await this.axiosInstance({
					url: options.url,
					method: options.method,
					data: options.requestData,
					...options.config
				});

				// Валидация ответа
				return options.responseSchema.parse(response.data);
			} catch (error) {
				if (error instanceof AxiosError) {
					// Обработка ошибок Axios
					if (error.response) {
						// Ошибка от сервера (4xx, 5xx)
						//throw new Error(`Server error: ${error.response.status} - ${error.message}`);
						console.log(error.response.data);
						return options.responseSchema.parse(error.response.data);
					} else if (error.request) {
						// Нет ответа от сервера
						throw new Error('No response received from server');
					} else {
						// Ошибка при настройке запроса
						throw new Error(`Request setup error: ${error.message}`);
					}
				} else if (error instanceof z.ZodError) {
					// Обработка ошибок валидации Zod
					throw new Error(`Validation error: ${error.issues.map((i) => i.message).join(', ')}`);
				} else {
					// Общая ошибка
					throw new Error(`Unexpected error: ${error}`);
				}
			}
		}
	}

	export const apiClient = new ApiClient('http://localhost:10015/rest');
</script>
