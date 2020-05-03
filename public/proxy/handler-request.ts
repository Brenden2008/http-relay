export default class HandlerRequest {
    constructor(private readonly response: Response) {
    }

    get url(): string {
        return this.headerValue('HttpRelay-Proxy-Url')
    }

    get method(): string {
        return this.headerValue('HttpRelay-Proxy-Method')
    }

    get scheme(): string {
        return this.headerValue('HttpRelay-Proxy-Scheme')
    }

    get host(): string {
        return this.headerValue('HttpRelay-Proxy-Host')
    }

    get path(): string {
        return this.headerValue('HttpRelay-Proxy-Path')
    }

    get query(): string | null {
        return this.response.headers.get('HttpRelay-Proxy-Query')
    }

    get queryParams(): URLSearchParams {
        return new URLSearchParams(this.query ?? '')
    }

    get fragment(): string | null {
        return this.response.headers.get('HttpRelay-Proxy-Fragment')
    }

    get headers(): Headers {
        return this.response.headers
    }

    get body(): ReadableStream<Uint8Array> | null {
        return this.response.body
    }

    public arrayBuffer(): Promise<ArrayBuffer> {
        return this.response.arrayBuffer()
    }

    public blob(): Promise<Blob> {
        return this.response.blob()
    }

    public formData(): Promise<FormData> {
        return this.response.formData()
    }

    public json(): Promise<JSON> {
        return this.response.json()
    }

    public text(): Promise<string> {
        return this.response.text()
    }

    public headerValue(name: string): string {
        let value = this.response.headers.get(name)
        if (!value) throw new Error(`Unable to find "${name}" header field.`)
        return value
    }
}