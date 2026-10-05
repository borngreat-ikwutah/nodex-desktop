import { IconCpu, IconLock, IconWaveSine } from "@tabler/icons-react";
import { ConsoleShell } from "./console-shell";

type VitalsCardProps = {
  platform: string;
  uptime: string;
  keyFingerprint: string;
};

function Vital({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode;
  label: string;
  value: string;
}) {
  return (
    <div className="flex items-center gap-3 rounded-2xl bg-white/[0.02] px-4 py-3 ring-1 ring-white/[0.06] ring-inset transition-all duration-700 ease-[cubic-bezier(0.32,0.72,0,1)] hover:bg-white/[0.05]">
      <span className="text-white/35">{icon}</span>
      <div className="min-w-0">
        <p className="text-[10px] tracking-[0.18em] text-white/30 uppercase">
          {label}
        </p>
        <p className="mt-0.5 truncate text-xs font-medium text-white/80">
          {value}
        </p>
      </div>
    </div>
  );
}

export function VitalsCard({
  platform,
  uptime,
  keyFingerprint,
}: VitalsCardProps) {
  return (
    <ConsoleShell
      eyebrow="Environment"
      title="Machine Vitals"
      icon={<IconCpu className="size-4" stroke={1.4} />}
      className="lg:col-span-5"
    >
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-1 xl:grid-cols-2">
        <Vital
          icon={<IconWaveSine className="size-4" stroke={1.4} />}
          label="Session Uptime"
          value={uptime}
        />
        <Vital
          icon={<IconCpu className="size-4" stroke={1.4} />}
          label="Platform"
          value={platform}
        />
        <Vital
          icon={<IconLock className="size-4" stroke={1.4} />}
          label="Seal"
          value="AES-256-GCM"
        />
        <Vital
          icon={<IconWaveSine className="size-4" stroke={1.4} />}
          label="Key"
          value={keyFingerprint || "—"}
        />
      </div>
    </ConsoleShell>
  );
}
