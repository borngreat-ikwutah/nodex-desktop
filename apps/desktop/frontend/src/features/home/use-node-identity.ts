import { useEffect, useState } from "react";
import { GetIdentity, GetNodeID } from "../../../wailsjs/go/main/App";

export type NodeIdentity = {
  node_id: string;
  uuid: string;
  public_key: string;
};

type IdentityState = {
  identity: NodeIdentity | null;
  nodeID: string | null;
  error: string | null;
  isLoading: boolean;
};

export function useNodeIdentity(): IdentityState {
  const [state, setState] = useState<IdentityState>({
    identity: null,
    nodeID: null,
    error: null,
    isLoading: true,
  });

  useEffect(() => {
    let active = true;

    Promise.all([GetNodeID(), GetIdentity()])
      .then(([nodeID, identity]) => {
        if (!active) return;
        setState({
          identity,
          nodeID,
          error: null,
          isLoading: false,
        });
      })
      .catch((cause: unknown) => {
        if (!active) return;
        setState({
          identity: null,
          nodeID: null,
          error: cause instanceof Error ? cause.message : String(cause),
          isLoading: false,
        });
      });

    return () => {
      active = false;
    };
  }, []);

  return state;
}
