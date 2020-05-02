/// <reference path="handler.ts" />
/// <reference path="route.ts" />

namespace HttpRelay.Proxy {
    class Routes {
        private readonly routes: Route[] = []

        constructor(
            private readonly abortSig: AbortSignal,
            private readonly wSecret: string
        ) {}

        public add(method: string, path: string, handlerFunc: HandlerFunc): void {
            let handler = new Handler(handlerFunc, this.abortSig, this.wSecret)
            let route = new Route(method, path, handler)
        }
    }
}