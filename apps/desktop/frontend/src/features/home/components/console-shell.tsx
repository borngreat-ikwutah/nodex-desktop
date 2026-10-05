import type { ReactNode } from "react";

type ConsoleShellProps = {
  eyebrow: string;
  title: string;
  icon: ReactNode;
  action?: ReactNode;
  className?: string;
  children: ReactNode;
};

/**
 * Double-bezel console card: an aluminium tray holding a recessed glass core.
 * Every major surface on the home screen nests through this shell.
 */
export function ConsoleShell({
  eyebrow,
  title,
  icon,
  action,
  className = "",
  children,
}: ConsoleShellProps) {
  return (
    <section
      className={`group relative rounded-[2rem] bg-white/[0.03] p-1.5 ring-1 ring-white/10 transition-all duration-700 ease-[cubic-bezier(0.32,0.72,0,1)] hover:bg-white/[0.05] hover:ring-white/20 ${className}`}
    >
      <div className="flex h-full flex-col overflow-hidden rounded-[calc(2rem-0.375rem)] bg-[#08080a] shadow-[inset_0_1px_1px_rgba(255,255,255,0.06),0_24px_60px_-30px_rgba(0,0,0,0.9)] ring-1 ring-white/[0.06] ring-inset">
        <header className="flex items-center justify-between gap-4 px-6 pt-6 pb-5">
          <div className="flex items-center gap-3">
            <span className="flex size-9 items-center justify-center rounded-full bg-white/[0.05] text-emerald-300/90 ring-1 ring-white/10 ring-inset transition-transform duration-700 ease-[cubic-bezier(0.32,0.72,0,1)] group-hover:scale-105">
              {icon}
            </span>
            <div>
              <p className="text-[10px] font-medium tracking-[0.22em] text-white/35 uppercase">
                {eyebrow}
              </p>
              <h2 className="mt-1 text-sm font-semibold tracking-tight text-white/90">
                {title}
              </h2>
            </div>
          </div>
          {action}
        </header>
        <div className="flex-1 px-6 pb-6">{children}</div>
      </div>
    </section>
  );
}
