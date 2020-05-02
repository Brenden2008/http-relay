// /// <reference path="handler-request.ts" />

namespace HttpRelay.Proxy {
    export type PlainHeaders = Headers | Record<string, string>
    export type RouteParams = string[]

    interface RespondParams {
        body?: any,
        status?: number
        headers?: PlainHeaders
        fileName?: string
        download?: boolean
    }

    export class HandlerFuncResult {
        constructor(
            public readonly body: any,
            public readonly status?: number,
            public readonly headers?: PlainHeaders,
            public readonly fileName?: string,
            public readonly download?: boolean
        ) {}
    }

    type Body = string | Blob | ArrayBuffer | FormData | URLSearchParams | ReadableStream | Promise<ArrayBuffer>

    function isBody(value: Body): value is Body {
        return typeof(value) === 'string'
            || value instanceof Blob
            || value instanceof ArrayBuffer
            || value instanceof FormData
            || value instanceof URLSearchParams
            || value instanceof ReadableStream
    }

    export class HandlerCtx {
        constructor(
            public readonly request: HandlerRequest,
            public readonly abortSig: AbortSignal,
            public readonly routeParams: RouteParams
        ) {}

        get serverId(): string {
            return this.request.headerValue('HttpRelay-Proxy-ServerId')
        }

        get jobId(): string {
            return this.request.headerValue('HttpRelay-Proxy-JobId')
        }

        public respond(result: RespondParams = {}): HandlerFuncResult {
            return new HandlerFuncResult(result.body, result.status, result.headers, result.fileName, result.download)
        }
    }
}