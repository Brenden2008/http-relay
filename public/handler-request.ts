namespace HttpRelay.Proxy {
    class HandlerRequest {
        constructor(private readonly response: Response) {
        }

        get url(): string | null {
            return this.response.headers.get('HttpRelay-Proxy-Url')
        }

        get method(): string | null {
            return this.response.headers.get('HttpRelay-Proxy-Method')
        }

        get scheme(): string | null {
            return this.response.headers.get('HttpRelay-Proxy-Scheme')
        }

        get host(): string | null {
            return this.response.headers.get('HttpRelay-Proxy-Host')
        }

        get path(): string | null {
            return this.response.headers.get('HttpRelay-Proxy-Path')
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
    }
}