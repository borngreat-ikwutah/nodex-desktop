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
