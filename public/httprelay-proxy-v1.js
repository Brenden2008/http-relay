export default class Httprelay {
    constructor(serverId, proxyUrl='https://staging.httprelay.io/proxy') {
        // WARNING!!! State is shared between parallel requests!
        this.serverId = serverId
        this.proxyUrl = proxyUrl
        this.url = `${this.proxyUrl}/${this.serverId}`
        this.routes = []
        this.errRetry = 0
        this.abortCtrl = new AbortController()
        this.abortSig = this.abortCtrl.signal
    }

    start(parallel=4) {
        if (typeof window !== 'undefined') window.addEventListener('beforeunload', () => this.stop())
        for (let i=0; i<parallel; i++) this.serve()
    }

    stop() {
        this.abortCtrl.abort()
    }

    serve(init=this.newCliReqInit()) {
        if (!this.abortSig.aborted) {
            fetch(this.url, init).then(resp => {
                console.log(resp.status)
                if (resp.status === 200) {
                    this.errRetry = 0
                    let method = resp.headers.get('httprelay-proxy-method')
                    let path = resp.headers.get('httprelay-proxy-path')
                    let handler = this.getHandler(method, path)
                    this.handle(resp, handler.func, handler.params)
                        .then(req => this.serve(req))
                } else {
                    this.handleError()
                }
            }, (err) => {
                console.log('Fetch err ', err)
                this.handleError(init)
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

    handleError(init=this.newCliReqInit()) {
        if (!this.abortSig.aborted) {
            setTimeout(() => this.serve(init), this.errRetry++ * 1000)
        }
    }

    getHandler(method, path) {
        let route = this.routes.find(r => method.match(r.method) && path.match(r.path))
        return route ? {
            func: route.handler,
            params: path.match(route.path).slice(1)
        } : {
            func: function() {
                return new Response(`Not handler for the "${method} ${path}" route on "${this.serverId}" server.`, { status: 404 })
            },
            params: []
        }
    }

    respToHandlerReq(resp) {
        return {
            method: resp.headers.get('httprelay-proxy-method'),
            url: resp.headers.get('httprelay-proxy-path'),
            headers: new Headers(resp.headers),
            body: resp.body,
            arrayBuffer: () => resp.arrayBuffer(),
            blob: () => resp.blob(),
            formData: () => resp.formData(),
            json: () => resp.json(),
            text: () => resp.text()
        }
    }

    respToCliReqInit(resp) {
        return resp.arrayBuffer().then(body => {
            let init =  this.newCliReqInit()
            init.headers = new Headers(resp.headers)
            init.headers.set('httprelay-proxy-status', resp.status || 200)
            init.body = body
            return init
        })
    }

    newCliReqInit() {
        return {
            method: 'SERVE',
            keepalive: true,
            signal: this.abortSig
        }
    }

    addRoute(method, path, handler) {
        this.routes.push({
            method: RegExp(method),
            path: new RegExp("^" + path.replace(/:[^\s/]+/g, '([\\w-]+)') + "$"),
            handler: handler
        })
        this.routes.sort((a, b) => b.path.length - a.path.length )
    }

    get(path, handler) {
        this.addRoute('GET', path, handler)
    }

    post(path, handler) {
        this.addRoute('POST', path, handler)
    }

    documentResponse(document) {
        return new Response(new XMLSerializer().serializeToString(document), {
            headers: { 'Content-Type': 'text/html; charset=UTF-8' }
        })
    }

    fileResponse(file, download = true) {
        return new Response(file, {
            headers: {
                'Httprelay-Proxy-Headers': 'Content-Disposition',
                'Content-Type': file.type,
                'Content-Disposition': `${download ? 'attachment' : 'inline'}; filename*=${this.encode(file.name)}`
            }
        })
    }

    encode(str) {
        return `UTF-8''` +
            encodeURIComponent(str)
                .replace(/['()]/g,
                    function(match) {
                        return '%' + match.charCodeAt(0).toString(16);
                    })
                .replace(/\*/g, '%2A')
                .replace(/%(7C|60|5E)/g,
                    function(_, match) {
                        return String.fromCharCode(parseInt(match, 16));
                    });
    }
}