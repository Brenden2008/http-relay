import HttprelayRouter from "./httprelay-proxy-router"
import HttprelayHandler from "./httprelay-proxy-handler"
import HttprelaySerResponse from './httprelay-proxy-ser-response'

export default class Httprelay {
    constructor(serverId=null, secret=null, proxyUrl='https://staging.httprelay.io/proxy') {
//    constructor(serverId, proxyUrl='http://localhost:8080/proxy') {
        // WARNING!!! State is shared between parallel requests!
        this._serverId = serverId || btoa(Math.random()).substr(5, 5)
        this._wSecret = secret || btoa(Math.random()).substr(8, 8)
        this._proxyUrl = proxyUrl
        this._url = `${this._proxyUrl}/${this._serverId}`
        this._errRetry = 0
        this._abortCtrl = new AbortController()
        this._abortSig = this._abortCtrl.signal
        this._routes = new HttprelayRouter(this._wSecret, this._abortSig)
    }

    start(parallel=4) {
        if (typeof window !== 'undefined') window.addEventListener('beforeunload', () => this.stop())
        for (let i=0; i<parallel; i++) this._serve()
    }

    stop() {
        this._abortCtrl.abort()
    }

    get routes() {
        return this._routes
    }

    _serve(initPro=null) {
        if (!this._abortSig.aborted) {
            initPro = initPro || new HttprelaySerResponse(this._wSecret, this._abortSig).reqInitPro
            Promise.resolve(initPro).then(init => {
                fetch(this._url, init).then(resp => {
                    if (resp.status === 200) {
                        this._errRetry = 0
                        let method = resp.headers.get('httprelay-proxy-method')
                        let path = resp.headers.get('httprelay-proxy-path')
                        let route = this._routes.getRoute(method, path)
                        CIA!!!!!!

                        this.handle(resp, route.handler, route.params)
                            .then(req => this._serve(req))
                    } else {
                        this.handleError(`Httprelay responded ${resp.status} while returning result and requesting new job`)
                    }
                }, err => this.handleError(err, init))
            })
        }
    }

    handle(resp, handlerFunc, handlerParams) {
        let respPro = handlerFunc(handlerParams, this.respToHandlerReq(resp))     // User can return promise or response
        return Promise.resolve(respPro)
            .then(r => this.respToCliReqInit(r))
            .then(req => {
                req.headers.set('httprelay-proxy-jobid', resp.headers.get('httprelay-proxy-jobid'))
                return req
            })
    }

    handleError(err, init=this.newCliReqInit()) {
        if (!this._abortSig.aborted) {
            setTimeout(() => this._serve(init), this._errRetry++ * 1000)
            throw err
        }
    }
}





// newCliReqInit(status = 200, headers = {}, body = null) {
//     let newHeaders = new Headers(headers)
//     newHeaders.set('httprelay-proxy-status', status)
//     if (this.wSecret) newHeaders.set('httprelay-wsecret', status)
//     return {
//         method: 'SERVE',
//         headers: newHeaders,
//         body: body,
//         signal: this._abortSig
//     }
// }
//
// respToCliReqInit(resp) {
//     switch (resp.constructor) {
//         case String:
//             return this.stringResponse(resp);
//         case Object:
//             return this.jsonResponse(resp);
//         case Array:
//             return this.jsonResponse(resp);
//         case Response:
//             return resp.arrayBuffer().then(body => this.newCliReqInit(resp.status, resp.headers, body))
//         default:
//             return resp
//     }
// }
//


// respToHandlerReq(resp) {
//     return {
//         method: resp.headers.get('httprelay-proxy-method'),
//         url: resp.headers.get('httprelay-proxy-path'),
//         headers: new Headers(resp.headers),
//         body: resp.body,
//         arrayBuffer: () => resp.arrayBuffer(),
//         blob: () => resp.blob(),
//         formData: () => resp.formData(),
//         json: () => resp.json(),
//         text: () => resp.text()
//     }
// }

// class HandlerResponse {
//     constructor(status = 200, headers = {}, body = null, abortSig) {
//         let newHeaders = new Headers(headers)
//         newHeaders.set('httprelay-proxy-status', status)
//         if (this.wSecret) newHeaders.set('httprelay-wsecret', status)
//         const this.init = {
//             method: 'SERVE',
//             headers: newHeaders,
//             body: body,
//             signal: abortSig
//         }
//     }
// }