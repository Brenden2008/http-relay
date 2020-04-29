import HttprelayHandler from "./httprelay-proxy-handler"
export default class HttprelayRouter {
    constructor(wSecret, abortSig) {
        this._routes = []
        this._wSecret = wSecret
        this._abortSig = abortSig
    }

    add(method, path, handlerFunc) {
        this._routes.push({
            method: RegExp(method),
            path: path,
            regExp: new RegExp("^" + path.replace(/:[^\s/]+/g, '([\\w-]+)') + "$"),
            handler: new HttprelayHandler(handlerFunc, this._wSecret, this._abortSig)
        })
        this._routes.sort((a, b) => b.regExp.length - a.regExp.length )
    }

    addGet(path, handler) {
        this.add('GET', path, handler)
    }

    addPost(path, handler) {
        this.add('POST', path, handler)
    }

    getRoute(method, path) {
        let route = this._routes.find(r => method.match(r.method) && path.match(r.regExp))
        return route ? {
            handler: route.handler,
            params: path.match(route.regExp).slice(1)
        } : {
            handler: new HttprelayHandler((req, params, resp) => {
                return resp.respond(`Not handler for the "${method} ${path}" route on "${this.serverId}" server.`, { status: 404 })
            }, this._wSecret, this._abortSig),
            params: []
        }
    }
}