import Routes from "./routes";
import Handler from "./handler";
import HandlerCtx from "./handler-ctx";
import HandlerRequest from "./handler-request";

export default class Proxy {
    public readonly routes: Routes;
    private readonly serverUrl
    private errRetry: number = 0;
    private readonly abortCtrl: AbortController;
    private readonly abortSig: AbortSignal;

    constructor(public readonly serverId: string, public readonly url: string, public readonly wSecret?: string) {
        // WARNING!!! State is shared between parallel requests!
        this.serverUrl = `${this.url}/${this.serverId}`
        this.routes = new Routes()
        this.abortCtrl = new AbortController()
        this.abortSig = this.abortCtrl.signal
    }

    public start(parallel=4): void {
        if (typeof window !== 'undefined') window.addEventListener('beforeunload', () => this.stop())
        for (let i=0; i<parallel; i++) this.serve()
    }

    public stop(): void {
        this.abortCtrl.abort()
    }

    private serve(init=Handler.requestInit(this.abortSig)) {
        if (!this.abortSig.aborted) {
            fetch(this.serverUrl, init).then(
                resp => {
                    if (resp.status === 200) {
                        this.errRetry = 0

                        let handlerRequest = new HandlerRequest(resp)
                        let selectedRoute = this.routes.find(handlerRequest.method, handlerRequest.path)
                        let ctx = new HandlerCtx(handlerRequest, this.abortSig, selectedRoute.params)
                        let handler = new Handler(selectedRoute.handlerFunc, this.wSecret)
                        handler.execute(ctx).then(
                            init => this.serve(init),
                            err => this.handleError(err)
                        )
                    } else {
                        this.handleError(`HttpRelay responded ${resp.status} while returning result and requesting new job.`)
                    }
                },
                err => this.handleError(err, init)
            )
        }
    }

    private handleError(err, init=Handler.requestInit(this.abortSig)) {
        if (!this.abortSig.aborted) {
            setTimeout(() => this.serve(init), this.errRetry++ * 1000)
            throw err
        }
    }
}