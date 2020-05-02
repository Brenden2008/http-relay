class HttpRelay {
    constructor(public readonly url: string='https://staging.httprelay.io') {
    }

    public proxy(serverId: string, wSecret?: string, path: string = '/proxy') {
        return new HRProxy(serverId, `${this.url}${path}`, wSecret)
    }

}