/// <reference path="handler.ts" />

namespace HttpRelay.Proxy {
    export class Route {
        private readonly methodRe: RegExp
        private readonly pathRe: RegExp
        private readonly pathDepth: number

        constructor(
            public readonly method: string,
            public readonly path: string,
            private readonly handler: Handler
        ) {
            this.methodRe = RegExp(method)
            this.pathRe = RegExp(path)
            this.pathDepth = this.path.split('/').length
        }

        public compare(r: Route): number {
            let result = r.pathDepth - this.pathDepth
            if (result == 0) result = r.path.length - this.path.length
            return result
        }
    }
}