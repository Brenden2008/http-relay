namespace HttpRelay.Proxy {
    type PlainHeaders = Headers | Record<string, string>

    interface ResponseMeta {
        status?: number
        headers?: PlainHeaders
        fileName?: string
        download?: boolean
    }

    type Body = string | Blob | ArrayBuffer | FormData | URLSearchParams | ReadableStream | Promise<ArrayBuffer>

    function isBody(value: Body): value is Body {
        return typeof(value) === "string" || value instanceof Blob || value instanceof ArrayBuffer || value instanceof FormData || value instanceof URLSearchParams || value instanceof ReadableStream
    }

    interface SerRespInit {
        method: string
        headers: Headers
        body: Body
        signal: AbortSignal
    }

    class HandlerCtx {
        public readonly request: HandlerRequest;

        constructor(private readonly cliResponse: Response, public readonly abortSig: AbortSignal, public readonly routeParams: string[]) {
            this.request = new HandlerRequest(cliResponse);
        }

        get serverId(): string | null {
            return this.cliResponse.headers.get('HttpRelay-Proxy-ServerId')
        }

        get jobId(): string | null {
            return this.cliResponse.headers.get('HttpRelay-Proxy-JobId')
        }

        public respond(result: any, meta: ResponseMeta = {}): Promise<SerRespInit> {
            return Promise.resolve(result)
                .then(r => this.getInitPro(r, meta.status, meta.headers, meta.fileName, meta.download))
        }

        private getInitPro(content: any, status?: number, customHeaders?: PlainHeaders, fileName?: string, download?: boolean): Promise<SerRespInit> {
            let headers: Headers
            let body: Body

            let defaultHeaders = new Headers()
            let defaultStatus = 200
            let defaultContentType = 'application/json'
            let defaultFileName = ''

            if (typeof content === 'string') { // STRING ////////////////////////////////////////////////////////////////////
                defaultContentType = 'text/html; charset=UTF-8'
                body = content
            } else if (content instanceof Document) { // DOCUMENT ///////////////////////////////////////////////////////////
                defaultContentType = 'text/html; charset=UTF-8'
                body = new XMLSerializer().serializeToString(content)
            } else if (content instanceof Response) { // RESPONSE ///////////////////////////////////////////////////////////
                defaultStatus = content.status
                defaultContentType = content.headers.get('content-type') ?? ''
                defaultHeaders = content.headers
                body = content.arrayBuffer()
            } else if (content instanceof File) { // FILE ///////////////////////////////////////////////////////////////////
                defaultContentType = content.type
                defaultFileName = content.name
                body = content
            } else if (isBody(content)) { // BODY ///////////////////////////////////////////////////////////////////
                body = content
            } else { // JSON ///////////////////////////////////////////////////////////////////////////////////////////////
                body = JSON.stringify(content)
            }

            headers = customHeaders ? new Headers(customHeaders) : defaultHeaders
            if (!headers.has('content-type')) headers.append('content-type', defaultContentType)
            if (fileName) defaultFileName = fileName
            if (download || defaultFileName) {
                let defaultContentDisposition = `${download ? 'attachment' : 'inline'};`
                if (fileName) defaultContentDisposition += ` filename*=${this.encode(fileName)}`
                if (!headers.has('content-disposition')) headers.append('content-disposition', defaultContentDisposition)
            }

            let headerWhitelist = Array.from(headers).map(h => h[0]).join(', ')
            headers.set('httprelay-proxy-headers', headerWhitelist) // Whitelisting headers that must be passed to client
            headers.set('httprelay-proxy-status', `${status ?? defaultStatus}`)

            return Promise.resolve(body)
                .then(b => (<SerRespInit> {
                    method: 'SERVE',
                    headers: headers,
                    body: b,
                    signal: this.abortSig
                }))
        }

        private static encode(str: string) {
            return `UTF-8''` + encodeURIComponent(str)
                .replace(/['()]/g,
                    function(match) {
                        return '%' + match.charCodeAt(0).toString(16);
                    })
                .replace(/\*/g, '%2A')
                .replace(/%(7C|60|5E)/g,
                    function(_, match) {
                        return String.fromCharCode(parseInt(match, 16));
                    });
        }

    }
}