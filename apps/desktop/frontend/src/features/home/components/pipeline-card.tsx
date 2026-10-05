import { IconAntennaBars5 } from "@tabler/icons-react";
import { ConsoleShell } from "./console-shell";

type PipelineCardProps = {
  nodeID: string | null;
  isLoading: boolean;
};

type Stage = {
  label: string;
  detail: string;
  state: "ready" | "standby" | "idle";
};

const stages: Stage[] = [
  { label: "Identity", detail: "Seeded, sealed, stable", state: "ready" },
  {
    label: "Local discovery",
    detail: "Awaiting Engine 02 build",
    state: "standby",
  },
  {
    label: "Transport link",
    detail: "TCP handshake not negotiated",
    state: "idle",
  },
];

const stateStyles: Record<Stage["state"], string> = {
  ready: "bg-emerald-400",
  standby: "bg-amber-300/70",
  idle: "bg-white/20",
};

export function PipelineCard({ nodeID, isLoading }: PipelineCardProps) {
  return (
    <ConsoleShell
      eyebrow="Agent Pipeline"
      title="Identity → Discovery → Link"
      icon={<IconAntennaBars5 className="size-4" stroke={1.4} />}
      className="lg:col-span-5"
    >
      <ul className="space-y-1">
        {stages.map((stage) => (
          <li
            key={stage.label}
            className="group/stage flex items-center gap-4 rounded-2xl px-3 py-3 transition-colors duration-500 ease-[cubic-bezier(0.32,0.72,0,1)] hover:bg-white/[0.03]"
          >
            <span className="relative flex size-2.5 shrink-0 items-center justify-center">
              <span
                className={`absolute size-2.5 rounded-full ${stateStyles[stage.state]} ${
                  stage.state === "ready" ? "animate-ping opacity-60" : ""
                }`}
              />
              <span
                className={`relative size-1.5 rounded-full ${stateStyles[stage.state]}`}
              />
            </span>
            <div className="min-w-0 flex-1">
              <p className="text-sm font-medium tracking-tight text-white/85">
                {stage.label}
              </p>
              <p className="mt-0.5 truncate text-xs text-white/35">
                {stage.detail}
              </p>
            </div>
            <span className="shrink-0 rounded-full bg-white/[0.04] px-2.5 py-1 font-mono text-[10px] tracking-[0.14em] text-white/45 uppercase ring-1 ring-white/10 ring-inset">
              {stage.state === "ready"
                ? "live"
                : isLoading && nodeID == null
                  ? "…"
                  : stage.state}
            </span>
          </li>
        ))}
      </ul>
    </ConsoleShell>
  );
}
