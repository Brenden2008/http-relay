// /// <reference path="handler-request.ts" />

namespace HttpRelay.Proxy {
    type PlainHeaders = Headers | Record<string, string>

    interface ResponseMeta {
        status?: number
        headers?: PlainHeaders
        fileName?: string
        download?: boolean
    }

    class HandlerResult {
        constructor(
            public readonly content: any,
            public readonly status?: number,
            public readonly headers?: PlainHeaders,
            public readonly fileName?: string,
            public readonly download?: boolean
        ) {}
    }

    type Body = string | Blob | ArrayBuffer | FormData | URLSearchParams | ReadableStream | Promise<ArrayBuffer>

    function isBody(value: Body): value is Body {
        return typeof(value) === "string"
            || value instanceof Blob
            || value instanceof ArrayBuffer
            || value instanceof FormData
            || value instanceof URLSearchParams
            || value instanceof ReadableStream
    }

    export class HandlerCtx {
        constructor(
            private readonly request: HandlerRequest,
            public readonly abortSig: AbortSignal,
            public readonly routeParams: string[]
        ) {}

        get serverId(): string {
            return this.request.headerValue('HttpRelay-Proxy-ServerId')
        }

        get jobId(): string {
            return this.request.headerValue('HttpRelay-Proxy-JobId')
        }

        public respond(content: any, meta: ResponseMeta = {}): HandlerResult {
            return new HandlerResult(content, meta.status, meta.headers, meta.fileName, meta.download)
        }
    }
}