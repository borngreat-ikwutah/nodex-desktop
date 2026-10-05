import { useState } from "react";
import { IconArrowUpRight, IconLink } from "@tabler/icons-react";
import { ConsoleShell } from "./console-shell";

/** Groups a 9 digit ID the way the engine renders it, tolerating partial input. */
function formatNodeIDInput(raw: string) {
  const digits = raw.replace(/\D/g, "").slice(0, 9);
  if (digits.length <= 3) return digits;
  if (digits.length <= 6) return `${digits.slice(0, 3)} ${digits.slice(3)}`;
  return `${digits.slice(0, 3)} ${digits.slice(3, 6)} ${digits.slice(6)}`;
}

export function ConnectCard() {
  const [value, setValue] = useState("");
  const digits = value.replace(/\D/g, "");
  const isComplete = digits.length === 9;

  return (
    <ConsoleShell
      eyebrow="Initiate"
      title="Reach a Remote Device"
      icon={<IconLink className="size-4" stroke={1.4} />}
      className="lg:col-span-7"
    >
      <form
        onSubmit={(event) => event.preventDefault()}
        className="flex flex-col gap-3 sm:flex-row sm:items-center"
      >
        <div className="flex flex-1 items-center rounded-full bg-white/[0.04] p-1.5 ring-1 ring-white/10 transition-all duration-500 ease-[cubic-bezier(0.32,0.72,0,1)] focus-within:bg-white/[0.07] focus-within:ring-emerald-300/40">
          <input
            value={value}
            onChange={(event) =>
              setValue(formatNodeIDInput(event.target.value))
            }
            inputMode="numeric"
            placeholder="000 000 000"
            aria-label="Remote Nodex ID"
            disabled={!isComplete}
            className="h-10 w-full min-w-0 rounded-full bg-transparent px-5 font-mono text-base tracking-[0.18em] text-white placeholder:text-white/20 focus:outline-none disabled:cursor-not-allowed"
          />
          <button
            type="submit"
            disabled={!isComplete}
            aria-label="Establish connection"
            className="group flex h-10 shrink-0 items-center gap-2 rounded-full bg-emerald-400 pr-1.5 pl-5 text-sm font-semibold text-emerald-950 transition-all duration-500 ease-[cubic-bezier(0.32,0.72,0,1)] hover:bg-emerald-300 active:scale-[0.98] disabled:bg-white/[0.06] disabled:text-white/25 disabled:pl-5"
          >
            Connect
            <span className="flex size-7 items-center justify-center rounded-full bg-emerald-950/15 transition-transform duration-500 ease-[cubic-bezier(0.32,0.72,0,1)] group-hover:translate-x-0.5 group-hover:-translate-y-px group-hover:scale-105">
              <IconArrowUpRight className="size-3.5" stroke={1.8} aria-hidden />
            </span>
          </button>
        </div>
      </form>
      <p className="mt-4 text-xs leading-relaxed text-white/35">
        Connect negotiates through the discovery engine, which is not built yet.
        A complete ID is required before a link can be requested.
      </p>
    </ConsoleShell>
  );
}
