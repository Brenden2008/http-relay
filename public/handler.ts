/// <reference path="handler-request.ts" />
/// <reference path="handler-ctx.ts" />

namespace HttpRelay.Proxy {
    export type HandlerFunc = (ctx: HandlerCtx) => any;

    interface SerRespInit {
        method: string
        headers: Headers
        body: Body
        signal: AbortSignal
    }

    export class Handler {
        constructor(
            private readonly handlerFunc: HandlerFunc,
            private readonly abortSig: AbortSignal,
            private readonly wSecret: string
        ) {}

// Add these before fetch
//    , private readonly abortSig: AbortSignal
//        if (wSecret) this.headers.set('httprelay-wsecret', wSecret)
//        signal: this.abortSig
        execute(request: HandlerRequest, routeParams: string[]) {
            let ctx = new HandlerCtx(request, this.abortSig, routeParams)
            this.handlerFunc(ctx)

        }

        public getInitPro(content: any, status?: number, customHeaders?: PlainHeaders, fileName?: string, download?: boolean): Promise<SerRespInit> {
            let headers: Headers
            let body: Body

            let defaultHeaders = new Headers()
            let defaultStatus = 200
            let defaultContentType: string
            let defaultFileName: string

            if (typeof content === 'undefined') { // EMPTY /////////////////////////////////////////////////////////////
                defaultStatus = 204
            } else if (typeof content === 'string') { // STRING ////////////////////////////////////////////////////////
                defaultContentType = 'text/html; charset=UTF-8'
                body = content
            } else if (content instanceof Document) { // DOCUMENT //////////////////////////////////////////////////////
                defaultContentType = 'text/html; charset=UTF-8'
                body = new XMLSerializer().serializeToString(content)
            } else if (content instanceof Response) { // RESPONSE //////////////////////////////////////////////////////
                defaultStatus = content.status
                defaultContentType = content.headers.get('content-type') ?? ''
                defaultHeaders = content.headers
                body = content.arrayBuffer()
            } else if (content instanceof File) { // FILE //////////////////////////////////////////////////////////////
                defaultContentType = content.type
                defaultFileName = content.name
                body = content
            } else if (content instanceof FormData) { // FORM DATA /////////////////////////////////////////////////////
                defaultContentType = 'multipart/form-data'
                body = content
            } else if (content instanceof URLSearchParams) { // URLENCODED /////////////////////////////////////////////
                defaultContentType = 'application/x-www-form-urlencoded; charset=utf-8'
                body = content
            } else if (isBody(content)) { // BODY //////////////////////////////////////////////////////////////////////
                defaultContentType = 'application/octet-stream'
                body = content
            } else { // JSON ///////////////////////////////////////////////////////////////////////////////////////////
                defaultContentType = 'application/json'
                body = JSON.stringify(content)
            }

            headers = customHeaders ? new Headers(customHeaders) : defaultHeaders
            if (!headers.has('content-type') && defaultContentType) headers.append('content-type', defaultContentType)
            if (fileName) defaultFileName = fileName
            if (download || defaultFileName) {
                let defaultContentDisposition = `${download ? 'attachment' : 'inline'};`
                if (defaultFileName) defaultContentDisposition += ` filename*=${this.encode(defaultFileName)}`
                if (!headers.has('content-disposition')) headers.append('content-disposition', defaultContentDisposition)
            }

            let headerWhitelist = Array.from(headers).map(h => h[0]).join(', ')
            headers.set('httprelay-proxy-headers', headerWhitelist) // Whitelisting headers that must be passed to client
            headers.set('httprelay-proxy-status', `${status ?? defaultStatus}`)
            if (this.wSecret) headers.set('httprelay-wsecret', this.wSecret)

            return Promise.resolve(body)
                .then(b => (<SerRespInit> {
                    method: 'SERVE',
                    headers: headers,
                    body: b,
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