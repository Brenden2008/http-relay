/// <reference path="handler.ts" />
/// <reference path="route.ts" />

namespace HttpRelay.Proxy {

    export interface SelectedRoute {
        handler: Handler,
        params: RouteParams
    }

    class Routes {
        private readonly routes: Route[] = []
        private readonly notFoundHandler: Handler

        constructor(
            private readonly wSecret: string,
            private readonly abortSig: AbortSignal,
            notFoundHandlerFunc?: HandlerFunc
        ) {
            if (notFoundHandlerFunc) {
                this.notFoundHandler = new Handler(this.wSecret, this.abortSig, notFoundHandlerFunc)
            } else {
                this.notFoundHandler = new Handler(this.wSecret, this.abortSig, ctx => ctx.respond({
                    status: 404,
                    body: `Handler not found for the "${ctx.request.method} ${ctx.request.path}" route on "${ctx.serverId}" server.`
                }))
            }
        }

        public add(method: string, path: string, handlerFunc: HandlerFunc): void {
            let handler = new Handler(this.wSecret, this.abortSig, handlerFunc)
            let route = new Route(method, path, handler)
            this.routes.push(route)
            this.routes.sort((a, b) => a.compare(b))
        }

        public find(method: string, path: string): SelectedRoute {
            let routeParams: RouteParams = []
            let route = this.routes.find(r => {
                let matchRes = r.match(method, path)
                if (matchRes != null) {
                    routeParams = matchRes
                    return true
                }
            })

            return <SelectedRoute> {
                handler: route ? route.handler : this.notFoundHandler,
                params: routeParams
            }
        }

        public get(path: string, handlerFunc: HandlerFunc): void {
            this.add('GET', path, handlerFunc)
        }

        public post(path: string, handlerFunc: HandlerFunc): void {
            this.add('POST', path, handlerFunc)
        }
    }
}