"use strict";
var HttpRelay;
(function (HttpRelay) {
    var Proxy;
    (function (Proxy) {
        class HandlerCtx {
            constructor(cliResponse, abortSig, routeParams) {
                this.cliResponse = cliResponse;
                this.abortSig = abortSig;
                this.routeParams = routeParams;
                this.request = new Proxy.HandlerRequest(cliResponse);
            }
            get serverId() {
                return this.cliResponse.headers.get('HttpRelay-Proxy-ServerId');
            }
            get jobId() {
                return this.cliResponse.headers.get('HttpRelay-Proxy-JobId');
            }
            respond(result, meta = {}) {
            }
        }
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
                return this.response.headers.get('HttpRelay-Proxy-Url');
            }
            get method() {
                return this.response.headers.get('HttpRelay-Proxy-Method');
            }
            get scheme() {
                return this.response.headers.get('HttpRelay-Proxy-Scheme');
            }
            get host() {
                return this.response.headers.get('HttpRelay-Proxy-Host');
            }
            get path() {
                return this.response.headers.get('HttpRelay-Proxy-Path');
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
        }
        Proxy.HandlerRequest = HandlerRequest;
    })(Proxy = HttpRelay.Proxy || (HttpRelay.Proxy = {}));
})(HttpRelay || (HttpRelay = {}));
var HttpRelay;
(function (HttpRelay) {
    var Proxy;
    (function (Proxy) {
        function isBody(value) {
            return typeof (value) === "string" || value instanceof Blob || value instanceof ArrayBuffer || value instanceof FormData || value instanceof URLSearchParams || value instanceof ReadableStream;
        }
        class HandlerResult {
            constructor(content, abortSig, wSecret, status, headers, fileName, download) {
                var _a;
                this.abortSig = abortSig;
                let defaultHeaders = new Headers();
                let defaultStatus = 200;
                let defaultContentType = 'application/json';
                let defaultFileName = '';
                if (typeof content === 'string') {
                    defaultContentType = 'text/html; charset=UTF-8';
                    this.body = content;
                }
                else if (content instanceof Document) {
                    defaultContentType = 'text/html; charset=UTF-8';
                    this.body = new XMLSerializer().serializeToString(content);
                }
                else if (content instanceof Response) {
                    defaultStatus = content.status;
                    defaultContentType = (_a = content.headers.get('content-type')) !== null && _a !== void 0 ? _a : '';
                    defaultHeaders = content.headers;
                    this.body = content.arrayBuffer();
                }
                else if (content instanceof File) {
                    defaultContentType = content.type;
                    defaultFileName = content.name;
                    this.body = content;
                }
                else if (isBody(content)) {
                    this.body = content;
                }
                else {
                    this.body = JSON.stringify(content);
                }
                this.headers = headers ? new Headers(headers) : defaultHeaders;
                if (!this.headers.has('content-type'))
                    this.headers.append('content-type', defaultContentType);
                if (fileName)
                    defaultFileName = fileName;
                if (download || defaultFileName) {
                    let defaultContentDisposition = `${download ? 'attachment' : 'inline'};`;
                    if (fileName)
                        defaultContentDisposition += ` filename*=${this.encode(fileName)}`;
                    if (!this.headers.has('content-disposition'))
                        this.headers.append('content-disposition', defaultContentDisposition);
                }
                let headerWhitelist = Array.from(this.headers).map(h => h[0]).join(', ');
                this.headers.set('httprelay-proxy-headers', headerWhitelist);
                this.headers.set('httprelay-proxy-status', `${status !== null && status !== void 0 ? status : defaultStatus}`);
                if (wSecret)
                    this.headers.set('httprelay-wsecret', wSecret);
            }
            get serRespInitPro() {
                return Promise.resolve(this.body)
                    .then(body => ({
                    method: 'SERVE',
                    headers: this.headers,
                    body: body,
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
    })(Proxy = HttpRelay.Proxy || (HttpRelay.Proxy = {}));
})(HttpRelay || (HttpRelay = {}));
var HttpRelay;
(function (HttpRelay) {
    var Proxy;
    (function (Proxy) {
        class Handler {
            constructor() {
            }
        }
    })(Proxy = HttpRelay.Proxy || (HttpRelay.Proxy = {}));
})(HttpRelay || (HttpRelay = {}));
//# sourceMappingURL=httprelay-proxy.js.map