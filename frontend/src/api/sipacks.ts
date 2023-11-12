import axios, { Axios, AxiosInstance } from "axios";

export class BackendRequester {
	api: AxiosInstance
	apiPrefix: string

	constructor(host: string, port: number, apiPrefix: string) {
		this.api = axios.create({
			baseURL: `http://${host}:${port}`,
		})
		this.apiPrefix = apiPrefix
	}

	async requestWithPrefix(method: string, url: string, data?: any) {
		return this.api.request({
			method,
			url: `${this.apiPrefix}${url}`,
			data,
		})
	}

	async getPacks() {
		return this.requestWithPrefix("GET", "/packs")
	}
}