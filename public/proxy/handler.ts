/// <reference path="handler-response.ts" />

//namespace HttpRelay.Proxy {
    class Handler {
        constructor(
            private readonly handlerFunc: HandlerFunc,
            private readonly wSecret?: string,
        ) {}

        public execute(ctx: HandlerCtx): Promise<RequestInit> {
            let rawResult = this.handlerFunc(ctx)
            let handlerResponse = (rawResult instanceof HandlerResponse) ? rawResult : new HandlerResponse(rawResult) // If result is raw any object, making it HandlerFuncResponse
            handlerResponse.headers.set('HttpRelay-Proxy-JobId', ctx.jobId)
            if (this.wSecret) handlerResponse.headers.set('HttpRelay-WSecret', this.wSecret)

            return Promise.resolve(handlerResponse.body)
                .then(b => Handler.requestInit(ctx.abortSig, b, handlerResponse.headers))
        }

        public static requestInit(abortSig: AbortSignal, body?: ResultBody, headers?: Headers) {
            return <RequestInit> {
                method: 'SERVE',
                headers: headers,
                body: body,
                signal: abortSig
            }
        }
    }
//}