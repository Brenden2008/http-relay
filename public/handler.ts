namespace HttpRelay.Proxy {
    export class Handler {
        constructor(
            private readonly handlerFunc: HandlerFunc,
            private readonly routeParams: RouteParams,
            private readonly wSecret: string,
            private readonly abortSig: AbortSignal,
        ) {}

        public execute(handlerRequest: HandlerRequest): Promise<RequestInit> {
            let ctx = new HandlerCtx(handlerRequest, this.abortSig, this.routeParams)
            let rawResult = this.handlerFunc(ctx)
            let handlerResponse = (rawResult instanceof HandlerResponse) ? rawResult : new HandlerResponse(rawResult) // If result is raw any object, making it HandlerFuncResponse
            handlerResponse.headers.set('HttpRelay-Proxy-JobId', ctx.jobId)
            if (this.wSecret) handlerResponse.headers.set('HttpRelay-WSecret', this.wSecret)

            return Promise.resolve(handlerResponse.body)
                .then(b => (<RequestInit> {
                    method: 'SERVE',
                    headers: handlerResponse.headers,
                    body: b,
                    signal: this.abortSig
                }))
        }
    }
}