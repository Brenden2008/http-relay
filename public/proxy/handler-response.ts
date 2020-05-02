//namespace HttpRelay.Proxy {
    type ResultBody = string | Blob | ArrayBuffer | FormData | URLSearchParams | ReadableStream | Promise<ArrayBuffer>

    function isResultBody(value: ResultBody): value is ResultBody {
        return typeof(value) === 'string'
            || value instanceof Blob
            || value instanceof ArrayBuffer
            || value instanceof FormData
            || value instanceof URLSearchParams
            || value instanceof ReadableStream
    }

    class HandlerResponse {
        public readonly body: ResultBody | undefined
        public readonly headers: Headers

        constructor(body: any, status?: number, headers?: PlainHeaders, fileName?: string, download?: boolean) {
            let defaultHeaders = new Headers()
            let defaultStatus = 200
            let defaultContentType: string | undefined
            let defaultFileName: string | undefined

            if (typeof body === 'undefined') { // EMPTY /////////////////////////////////////////////////////////////
                defaultStatus = 204
            } else if (typeof body === 'string') { // STRING ////////////////////////////////////////////////////////
                defaultContentType = 'text/html; charset=UTF-8'
                this.body = body
            } else if (body instanceof Document) { // DOCUMENT //////////////////////////////////////////////////////
                defaultContentType = 'text/html; charset=UTF-8'
                this.body = new XMLSerializer().serializeToString(body)
            } else if (body instanceof Response) { // RESPONSE //////////////////////////////////////////////////////
                defaultStatus = body.status
                defaultContentType = body.headers.get('content-type') ?? ''
                defaultHeaders = body.headers
                this.body = body.arrayBuffer()
            } else if (body instanceof File) { // FILE //////////////////////////////////////////////////////////////
                defaultContentType = body.type
                defaultFileName = body.name
                this.body = body
            } else if (body instanceof FormData) { // FORM DATA /////////////////////////////////////////////////////
                defaultContentType = 'multipart/form-data'
                this.body = body
            } else if (body instanceof URLSearchParams) { // URLENCODED /////////////////////////////////////////////
                defaultContentType = 'application/x-www-form-urlencoded; charset=utf-8'
                this.body = body
            } else if (isResultBody(body)) { // BODY ///////////////////////////////////////////////////////////////////
                defaultContentType = 'application/octet-stream'
                this.body = body
            } else { // JSON ///////////////////////////////////////////////////////////////////////////////////////////
                defaultContentType = 'application/json'
                this.body = JSON.stringify(body)
            }

            this.headers = headers ? new Headers(headers) : defaultHeaders
            if (!this.headers.has('content-type') && defaultContentType) this.headers.append('content-type', defaultContentType)
            if (fileName) defaultFileName = fileName
            if (download || defaultFileName) {
                let defaultContentDisposition = `${download ? 'attachment' : 'inline'};`
                if (defaultFileName) defaultContentDisposition += ` filename*=${this.encode(defaultFileName)}`
                if (!this.headers.has('content-disposition')) this.headers.append('content-disposition', defaultContentDisposition)
            }

            let headerWhitelist = Array.from(this.headers).map(h => h[0]).join(', ')
            this.headers.set('httprelay-proxy-headers', headerWhitelist) // Whitelisting headers that must be passed to client
            this.headers.set('httprelay-proxy-status', `${status ?? defaultStatus}`)
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
//}