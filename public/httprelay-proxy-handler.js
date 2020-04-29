import SerResponse from './httprelay-proxy-ser-response'

export default class HttprelayHandler {
    constructor(func, wSecret, abortSig) {
        this._handlerFunc = func
        this._wSecret = wSecret
        this._abortSig = abortSig
    }

    execute(resp, params) {
        let jonId = resp.headers.get('httprelay-proxy-jobid')
        let serResp = new SerResponse(this._wSecret, this._abortSig)
        let respPro = this._handlerFunc(this._respToHandlerReq(resp, params, serResp))     // User can return promise or response

        return Promise.resolve(respPro)
            .then(r => r instanceof hrResponse ? r : new hrResponse(r) )
            .then(r => {
            })
            .then(req => {
                req.headers.set('httprelay-proxy-jobid', resp.headers.get('httprelay-proxy-jobid'))
                return req
            })
    }

    _respToHandlerReq(resp, params) {
        return {
            method: resp.headers.get('httprelay-proxy-method'),
            url: resp.headers.get('httprelay-proxy-path'),
            headers: new Headers(resp.headers),
            params: params,
            body: resp.body,
            arrayBuffer: () => resp.arrayBuffer(),
            blob: () => resp.blob(),
            formData: () => resp.formData(),
            json: () => resp.json(),
            text: () => resp.text()
        }
    }

    _newCliReqInit(status, headers, body) {
        let newHeaders = new Headers(headers)
        newHeaders.set('httprelay-proxy-status', status)
        if (this.wSecret) newHeaders.set('httprelay-wsecret', status)
        return {
            method: 'SERVE',
            headers: newHeaders,
            body: body,
            signal: this.abortSig
        }
    }
}




















// _stringResponse(status, headers, string) {
//     let newHeaders = new Headers(headers)
//     newHeaders.set('content-type', 'text/html; charset=UTF-8')
//     return this._newCliReqInit(status, newHeaders, string)
// }
//
// _jsonResponse(status, headers, obj) {
//     let newHeaders = new Headers(headers)
//     newHeaders.set('content-type', 'application/json')
//     return this._newCliReqInit(status, newHeaders, JSON.stringify(obj))
// }
//
// _documentResponse(status, headers, document}) {
//     let newHeaders = new Headers(headers)
//     newHeaders.set('content-type', 'text/html; charset=UTF-8')
//     return this._newCliReqInit(status, newHeaders, new XMLSerializer().serializeToString(document))
// }
//
// _fileResponse(file, download = true, status = 200, headers = {}) {
//     let newHeaders = new Headers(headers)
//     newHeaders.set('content-type', file.type)
//     newHeaders.set('content-disposition', `${download ? 'attachment' : 'inline'}; filename*=${this.encode(file.name)}`)
//     newHeaders.set('httprelay-proxy-headers', 'content-disposition')
//     return this._newCliReqInit(status, newHeaders, file)
// }



// //--------
//
//
// this.headers.append('content-disposition', val)
//
// switch (result.constructor) {
//     case String:
//         defaultContentType = 'text/html; charset=UTF-8'
//         break
//     case Document:
//         let newHeaders = new Headers(headers)
//         newHeaders.set('content-type', 'text/html; charset=UTF-8')
//         return this._newCliReqInit(status, newHeaders, new XMLSerializer().serializeToString(document))
//
//         return this._newCliReqInit(resp.status, resp.headers, resp.body)
//     case Response:
//         defaultStatus = result.status
//         defaultContentType = result.headers.get('content-type')
//         defaultHeaders = result.headers
//         this.body = result.arrayBuffer()
//         break
//     default:
//         return this._jsonResponse(resp)
// }
//
// this.status = meta.status || 200
// this.headers = meta.headers || new Headers()
// this.contentType = meta.contentType || null
//
// if (!this.headers.has('content-disposition')) {
//     if (this.download || this.fileName) {
//         let val = `${this.download ? 'attachment' : 'inline'};`
//         if (this.fileName) val += ` filename*=${this.encode(this.fileName)}`
//         this.headers.append('content-disposition', val)
//         newHeaders.set('httprelay-proxy-headers', 'content-disposition')
//     }
// }
//
//
