import HandlerRequest from "./handler-request";
import HandlerResponse from "./handler-response";

export type PlainHeaders = Headers | Record<string, string>
export type RouteParams = string[]

interface RespondParams {
    body?: any,
    status?: number
    headers?: PlainHeaders
    fileName?: string
    download?: boolean
}

export default class HandlerCtx {
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

    public respond(result: RespondParams = {}): HandlerResponse {
        return new HandlerResponse(result.body, result.status, result.headers, result.fileName, result.download)
    }
}