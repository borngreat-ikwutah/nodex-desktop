export namespace discovery {
  export class Peer {
    node_id: string;
    uuid: string;
    public_key: string;
    ips: string[];
    port: number;

    static createFrom(source: any = {}) {
      return new Peer(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.node_id = source["node_id"];
      this.uuid = source["uuid"];
      this.public_key = source["public_key"];
      this.ips = source["ips"];
      this.port = source["port"];
    }
  }
}

export namespace identity {
  export class Public {
    node_id: string;
    uuid: string;
    public_key: string;

    static createFrom(source: any = {}) {
      return new Public(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.node_id = source["node_id"];
      this.uuid = source["uuid"];
      this.public_key = source["public_key"];
    }
  }
}

export namespace media {
  export class Signal {
    session_id: string;
    sdp?: string;
    candidate?: string;
    sdp_mline_index?: number;

    static createFrom(source: any = {}) {
      return new Signal(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.session_id = source["session_id"];
      this.sdp = source["sdp"];
      this.candidate = source["candidate"];
      this.sdp_mline_index = source["sdp_mline_index"];
    }
  }
}

export namespace transport {
  export class PingResult {
    peer_id: string;
    peer_uuid: string;
    peer_public_key: string;
    latency: number;
    latency_ms: number;
    remote_addr: string;

    static createFrom(source: any = {}) {
      return new PingResult(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.peer_id = source["peer_id"];
      this.peer_uuid = source["peer_uuid"];
      this.peer_public_key = source["peer_public_key"];
      this.latency = source["latency"];
      this.latency_ms = source["latency_ms"];
      this.remote_addr = source["remote_addr"];
    }
  }
}
