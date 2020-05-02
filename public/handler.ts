/// <reference path="handler-request.ts" />
/// <reference path="handler-ctx.ts" />

namespace HttpRelay.Proxy {
    export type HandlerFunc = (ctx: HandlerCtx) => any;

    export class Handler {
        constructor(
            private readonly wSecret: string,
            private readonly abortSig: AbortSignal,
            private readonly handlerFunc: HandlerFunc
        ) {}

        execute(request: HandlerRequest, routeParams: string[]): any {
            let ctx = new HandlerCtx(request, this.abortSig, routeParams)
            return this.handlerFunc(ctx)
        }
    }
}