/// <reference path="handler.ts" />

namespace HttpRelay.Proxy {
    class Route {
        private readonly methodRe: RegExp
        private readonly pathRe: RegExp

        constructor(public readonly method: string, public readonly path: string, private readonly handler: Handler) {
            this.methodRe = RegExp(method)
            this.pathRe = RegExp(path)
        }

        public toString() {

        }
    }
}