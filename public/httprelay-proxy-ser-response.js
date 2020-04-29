export default class SerResponse {
    constructor(wSecret, abortSig) {
        this._wSecret = wSecret
        this._abortSig = abortSig
    }

    respond(result, meta={}) {
        this._body = result

        let defaultHeaders = new Headers()
        let defaultStatus = 200
        let defaultContentType = 'application/json'
        let defaultContentDisposition = null
        let defaultFileName = null

        if (typeof result === 'string') { // STRING ////////////////////////////////////////////////////////////////////
            defaultContentType = 'text/html; charset=UTF-8'
        } else if (result instanceof Document) { // DOCUMENT ///////////////////////////////////////////////////////////
            defaultContentType = 'text/html; charset=UTF-8'
            this._body = new XMLSerializer().serializeToString(result)
        } else if (result instanceof Response) { // RESPONSE ///////////////////////////////////////////////////////////
            defaultStatus = result.status
            defaultContentType = result.headers.get('content-type')
            defaultHeaders = result.headers
            this._body = result.arrayBuffer()
        } else if (result instanceof File) { // FILE ///////////////////////////////////////////////////////////////////
            defaultContentType = result.type
            defaultFileName = result.name
        } else { // JSON ///////////////////////////////////////////////////////////////////////////////////////////////
            this._body = JSON.stringify(result)
        }

        if (meta.fileName) defaultFileName = meta.fileName
        if (meta.download || defaultFileName) {
            defaultContentDisposition = `${meta.download ? 'attachment' : 'inline'};`
            if (meta.fileName) defaultContentDisposition += ` filename*=${this._encode(meta.fileName)}`
        }

        this._headers = meta.headers || defaultHeaders
        if (!this._headers.has('content-type')) this._headers.append('content-type', defaultContentType)
        if (!this._headers.has('content-disposition')) this._headers.append('content-disposition', defaultContentDisposition)

        // Whitelisting headers that must be passed to client
        let headerWhitelist = Array.from(this._headers).map(h => h[0]).join(', ')
        this._headers.set('httprelay-proxy-headers', headerWhitelist)

        this._headers.set('httprelay-proxy-status', meta.status || defaultStatus)
        if (this._wSecret) this.headers.set('httprelay-wsecret', this._wSecret)
    }

    _encode(str) {
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

    get reqInitPro() {
        return Promise.resolve(this._body)
            .then(body => ({
                method: 'SERVE',
                headers: this._headers,
                body: body,
                signal: this._abortSig
            }))
    }

    get abortSig() {
        return this._abortSig
    }
}




