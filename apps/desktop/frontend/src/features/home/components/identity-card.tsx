import { useState } from "react";
import { toast } from "sonner";
import {
  IconCheck,
  IconCopy,
  IconFingerprint,
  IconShieldCheck,
} from "@tabler/icons-react";
import type { NodeIdentity } from "../use-node-identity";
import { ConsoleShell } from "./console-shell";

type IdentityCardProps = {
  nodeID: string | null;
  identity: NodeIdentity | null;
  isLoading: boolean;
  error: string | null;
};

function CopyButton({ value, label }: { value: string; label: string }) {
  const [copied, setCopied] = useState(false);

  async function handleCopy() {
    try {
      await navigator.clipboard.writeText(value);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 1600);
      if (Math.random() > 0.5) {
        toast("password has been copied, make it seen or unseen");
      }
    } catch {
      setCopied(false);
    }
  }

  return (
    <button
      type="button"
      onClick={handleCopy}
      aria-label={`Copy ${label}`}
      className="flex size-8 shrink-0 items-center justify-center rounded-full bg-white/[0.06] text-white/60 ring-1 ring-white/10 transition-all duration-500 ease-[cubic-bezier(0.32,0.72,0,1)] hover:bg-white/[0.12] hover:text-white active:scale-95"
    >
      {copied ? (
        <IconCheck className="size-3.5 text-emerald-300" />
      ) : (
        <IconCopy className="size-3.5" />
      )}
    </button>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-4 border-t border-white/[0.06] py-3">
      <span className="text-[11px] tracking-[0.16em] text-white/35 uppercase">
        {label}
      </span>
      <span className="truncate font-mono text-xs text-white/70">{value}</span>
    </div>
  );
}

export function IdentityCard({
  nodeID,
  identity,
  isLoading,
  error,
}: IdentityCardProps) {
  const fingerprint = identity?.public_key?.slice(0, 24) ?? "";
  const displayID = isLoading ? "000 000 000" : (nodeID ?? "unavailable");

  return (
    <ConsoleShell
      eyebrow="Identity Engine"
      title="Persistent Machine Identity"
      icon={<IconFingerprint className="size-4" stroke={1.4} />}
      action={
        <span className="flex items-center gap-1.5 rounded-full bg-emerald-400/10 px-3 py-1.5 text-[10px] font-medium tracking-[0.18em] text-emerald-300 uppercase ring-1 ring-emerald-400/20 ring-inset">
          <IconShieldCheck className="size-3" stroke={1.6} />
          Sealed
        </span>
      }
      className="lg:col-span-7 lg:row-span-2"
    >
      <div className="flex h-full flex-col justify-between gap-8">
        <div>
          <p className="text-[10px] tracking-[0.22em] text-white/30 uppercase">
            Your Nodex ID
          </p>
          <div className="mt-4 flex items-center gap-4">
            <p
              aria-live="polite"
              className="bg-gradient-to-br from-white via-white/90 to-emerald-200/70 bg-clip-text font-mono text-[clamp(2.5rem,7vw,4.25rem)] leading-none font-semibold tracking-[0.04em] text-transparent tabular-nums"
            >
              {displayID}
            </p>
            {nodeID ? <CopyButton value={nodeID} label="Nodex ID" /> : null}
          </div>
          <p className="mt-4 max-w-md text-sm leading-relaxed text-white/40">
            Minted once from 32 bytes of system entropy and sealed against this
            machine. Restart, reboot, or change networks: the ID never moves.
          </p>
        </div>

        <div className="rounded-[1.25rem] bg-black/40 p-1 ring-1 ring-white/[0.06] ring-inset">
          <div className="rounded-[calc(1.25rem-0.25rem)] bg-white/[0.02] px-4 py-1">
            <Row label="Installation UUID" value={identity?.uuid ?? "—"} />
            <Row
              label="Public Fingerprint"
              value={fingerprint ? `${fingerprint}…` : "—"}
            />
            <Row
              label="Key File"
              value="~/.config/Nodex/identity/machine.key"
            />
            <Row
              label="Integrity"
              value={error ? "unreadable" : isLoading ? "reading" : "verified"}
            />
          </div>
        </div>
      </div>
    </ConsoleShell>
  );
}
