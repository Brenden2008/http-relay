"use strict";
var HttpRelay;
(function (HttpRelay) {
    var Proxy;
    (function (Proxy) {
        function isBody(value) {
            return typeof (value) === "string" || value instanceof Blob || value instanceof ArrayBuffer || value instanceof FormData || value instanceof URLSearchParams || value instanceof ReadableStream;
        }
        class HandlerCtx {
            constructor(request, abortSig, routeParams) {
                this.request = request;
                this.abortSig = abortSig;
                this.routeParams = routeParams;
            }
            get serverId() {
                return this.request.headerValue('HttpRelay-Proxy-ServerId');
            }
            get jobId() {
                return this.request.headerValue('HttpRelay-Proxy-JobId');
            }
            respond(result, meta = {}) {
                return Promise.resolve(result)
                    .then(r => this.getInitPro(r, meta.status, meta.headers, meta.fileName, meta.download));
            }
            getInitPro(content, status, customHeaders, fileName, download) {
                var _a;
                let headers;
                let body;
                let defaultHeaders = new Headers();
                let defaultStatus = 200;
                let defaultContentType = 'application/json';
                let defaultFileName = '';
                if (typeof content === 'string') {
                    defaultContentType = 'text/html; charset=UTF-8';
                    body = content;
                }
                else if (content instanceof Document) {
                    defaultContentType = 'text/html; charset=UTF-8';
                    body = new XMLSerializer().serializeToString(content);
                }
                else if (content instanceof Response) {
                    defaultStatus = content.status;
                    defaultContentType = (_a = content.headers.get('content-type')) !== null && _a !== void 0 ? _a : '';
                    defaultHeaders = content.headers;
                    body = content.arrayBuffer();
                }
                else if (content instanceof File) {
                    defaultContentType = content.type;
                    defaultFileName = content.name;
                    body = content;
                }
                else if (isBody(content)) {
                    body = content;
                }
                else {
                    body = JSON.stringify(content);
                }
                headers = customHeaders ? new Headers(customHeaders) : defaultHeaders;
                if (!headers.has('content-type'))
                    headers.append('content-type', defaultContentType);
                if (fileName)
                    defaultFileName = fileName;
                if (download || defaultFileName) {
                    let defaultContentDisposition = `${download ? 'attachment' : 'inline'};`;
                    if (fileName)
                        defaultContentDisposition += ` filename*=${this.encode(fileName)}`;
                    if (!headers.has('content-disposition'))
                        headers.append('content-disposition', defaultContentDisposition);
                }
                let headerWhitelist = Array.from(headers).map(h => h[0]).join(', ');
                headers.set('httprelay-proxy-headers', headerWhitelist);
                headers.set('httprelay-proxy-status', `${status !== null && status !== void 0 ? status : defaultStatus}`);
                return Promise.resolve(body)
                    .then(b => ({
                    method: 'SERVE',
                    headers: headers,
                    body: b,
                    signal: this.abortSig
                }));
            }
            encode(str) {
                return `UTF-8''` + encodeURIComponent(str)
                    .replace(/['()]/g, function (match) {
                    return '%' + match.charCodeAt(0).toString(16);
                })
                    .replace(/\*/g, '%2A')
                    .replace(/%(7C|60|5E)/g, function (_, match) {
                    return String.fromCharCode(parseInt(match, 16));
                });
            }
        }
        Proxy.HandlerCtx = HandlerCtx;
    })(Proxy = HttpRelay.Proxy || (HttpRelay.Proxy = {}));
})(HttpRelay || (HttpRelay = {}));
var HttpRelay;
(function (HttpRelay) {
    var Proxy;
    (function (Proxy) {
        class HandlerRequest {
            constructor(response) {
                this.response = response;
            }
            get url() {
                return this.headerValue('HttpRelay-Proxy-Url');
            }
            get method() {
                return this.headerValue('HttpRelay-Proxy-Method');
            }
            get scheme() {
                return this.headerValue('HttpRelay-Proxy-Scheme');
            }
            get host() {
                return this.headerValue('HttpRelay-Proxy-Host');
            }
            get path() {
                return this.headerValue('HttpRelay-Proxy-Path');
            }
            get query() {
                return this.response.headers.get('HttpRelay-Proxy-Query');
            }
            get queryParams() {
                var _a;
                return new URLSearchParams((_a = this.query) !== null && _a !== void 0 ? _a : '');
            }
            get fragment() {
                return this.response.headers.get('HttpRelay-Proxy-Fragment');
            }
            get headers() {
                return this.response.headers;
            }
            get body() {
                return this.response.body;
            }
            arrayBuffer() {
                return this.response.arrayBuffer();
            }
            blob() {
                return this.response.blob();
            }
            formData() {
                return this.response.formData();
            }
            json() {
                return this.response.json();
            }
            text() {
                return this.response.text();
            }
            headerValue(name) {
                let value = this.response.headers.get(name);
                if (!value)
                    throw new Error(`Unable to find "${name}" header field.`);
                return value;
            }
        }
        Proxy.HandlerRequest = HandlerRequest;
    })(Proxy = HttpRelay.Proxy || (HttpRelay.Proxy = {}));
})(HttpRelay || (HttpRelay = {}));
var HttpRelay;
(function (HttpRelay) {
    var Proxy;
    (function (Proxy) {
        class Handler {
            constructor(handlerFunc, abortSig) {
                this.handlerFunc = handlerFunc;
                this.abortSig = abortSig;
            }
            execute(request, routeParams) {
                let ctx = new Proxy.HandlerCtx(request, this.abortSig, routeParams);
                this.handlerFunc(ctx);
            }
        }
        Proxy.Handler = Handler;
    })(Proxy = HttpRelay.Proxy || (HttpRelay.Proxy = {}));
})(HttpRelay || (HttpRelay = {}));
var HttpRelay;
(function (HttpRelay) {
    var Proxy;
    (function (Proxy) {
        class Route {
            constructor(method, path, handler) {
                this.method = method;
                this.path = path;
                this.handler = handler;
                this.methodRe = RegExp(method);
                this.pathRe = RegExp(path);
            }
            toString() {
            }
        }
    })(Proxy = HttpRelay.Proxy || (HttpRelay.Proxy = {}));
})(HttpRelay || (HttpRelay = {}));
//# sourceMappingURL=httprelay-proxy.js.map