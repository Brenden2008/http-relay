namespace HttpRelay.Proxy {
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
    
    class HandlerResult {
        private readonly headers: Headers
        private readonly body: Body

        constructor(content: any, private readonly abortSig: AbortSignal, wSecret: string, status?: number, headers?: Headers | Record<string, string>, fileName?: string, download?: boolean) {
            let defaultHeaders = new Headers()
            let defaultStatus = 200
            let defaultContentType = 'application/json'
            let defaultFileName = ''

            if (typeof content === 'string') { // STRING ////////////////////////////////////////////////////////////////////
                defaultContentType = 'text/html; charset=UTF-8'
                this.body = content
            } else if (content instanceof Document) { // DOCUMENT ///////////////////////////////////////////////////////////
                defaultContentType = 'text/html; charset=UTF-8'
                this.body = new XMLSerializer().serializeToString(content)
            } else if (content instanceof Response) { // RESPONSE ///////////////////////////////////////////////////////////
                defaultStatus = content.status
                defaultContentType = content.headers.get('content-type') ?? ''
                defaultHeaders = content.headers
                this.body = content.arrayBuffer()
            } else if (content instanceof File) { // FILE ///////////////////////////////////////////////////////////////////
                defaultContentType = content.type
                defaultFileName = content.name
                this.body = content
            } else if (isBody(content)) { // BODY ///////////////////////////////////////////////////////////////////
                this.body = content
            } else { // JSON ///////////////////////////////////////////////////////////////////////////////////////////////
                this.body = JSON.stringify(content)
            }

            this.headers = headers ? new Headers(headers) : defaultHeaders
            if (!this.headers.has('content-type')) this.headers.append('content-type', defaultContentType)
            if (fileName) defaultFileName = fileName
            if (download || defaultFileName) {
                let defaultContentDisposition = `${download ? 'attachment' : 'inline'};`
                if (fileName) defaultContentDisposition += ` filename*=${this.encode(fileName)}`
                if (!this.headers.has('content-disposition')) this.headers.append('content-disposition', defaultContentDisposition)
            }

            let headerWhitelist = Array.from(this.headers).map(h => h[0]).join(', ')
            this.headers.set('httprelay-proxy-headers', headerWhitelist) // Whitelisting headers that must be passed to client
            this.headers.set('httprelay-proxy-status', `${status ?? defaultStatus}`)
            if (wSecret) this.headers.set('httprelay-wsecret', wSecret)
        }

        get serRespInitPro(): Promise<SerRespInit> {
            return Promise.resolve(this.body)
                .then(body => (<SerRespInit> {
                    method: 'SERVE',
                    headers: this.headers,
                    body: body,
                    signal: this.abortSig
                }))
        }

        private encode(str: string) {
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