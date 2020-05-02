/// <reference path="handler.ts" />

namespace HttpRelay.Proxy {
    export class Route {
        private readonly methodRx: RegExp
        private readonly pathRx: RegExp
        private readonly pathDepth: number

        constructor(
            public readonly method: string,
            public readonly path: string,
            public readonly handler: Handler
        ) {
            this.methodRx = RegExp(method)
            this.pathRx = RegExp("^" + path.replace(/:[^\s/]+/g, '([\\w-]+)') + "$")
            this.pathDepth = this.path.split('/').length
        }

        public compare(r: Route): number {
            let result = r.pathDepth - this.pathDepth
            if (result == 0) result = r.path.length - this.path.length
            return result
        }

        public match(method: string, path: string): RouteParams | null {
            if (method.match(this.methodRx)) {
                let routeParams = path.match(this.pathRx)
                if (routeParams) return routeParams.slice(1)
            }
            return null
        }
    }
}