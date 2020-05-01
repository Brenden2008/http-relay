namespace HttpRelay.Proxy {
    interface ResponseMeta {
        status?: number
        headers?: Headers | Object
        fileName?: string
        download?: boolean
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

        public respond(result: any, meta: ResponseMeta = {}) {

        }
    }
}