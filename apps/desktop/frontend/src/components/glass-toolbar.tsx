import { useEffect, useRef } from "react";
import { useNavigate } from "@tanstack/react-router";
import {
  IconCircleFilled,
  IconMoon,
  IconSearch,
  IconSun,
} from "@tabler/icons-react";
import { Button } from "@repo/ui/components/ui/button.tsx";
import { Kbd } from "@repo/ui/components/ui/kbd.tsx";
import { useTheme } from "../lib/theme";
import { cn } from "../lib/cn";

type GlassToolbarProps = {
  status?: "online" | "offline";
  onConnect?: () => void;
  onSearch?: (query: string) => void;
};

/**
 * Floating glass toolbar (NSToolbar analogue).
 * Sits in the native traffic-light drag region; controls opt out of drag.
 */
export function GlassToolbar({
  status = "offline",
  onConnect,
  onSearch,
}: GlassToolbarProps) {
  const { mode, resolved, setMode } = useTheme();
  const navigate = useNavigate();
  const searchRef = useRef<HTMLInputElement>(null);
  const online = status === "online";

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      const mod = e.metaKey || e.ctrlKey;
      if (mod && e.key.toLowerCase() === "k") {
        e.preventDefault();
        searchRef.current?.focus();
      }
      if (mod && ["1", "2", "3", "4", "5"].includes(e.key)) {
        e.preventDefault();
        const targets = [
          "/home",
          "/devices",
          "/sessions",
          "/transfers",
          "/settings",
        ];
        navigate({ to: targets[Number(e.key) - 1] as never });
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [navigate]);

  const toggleTheme = () => {
    if (mode === "system") setMode(resolved === "dark" ? "light" : "dark");
    else if (mode === "dark") setMode("light");
    else setMode("dark");
  };

  return (
    <div className="mac-drag flex shrink-0 items-center gap-2 px-3 pt-3 sm:gap-3 sm:px-4">
      {/* Traffic-light safe inset */}
      <div className="w-[var(--mac-traffic-inset)] shrink-0" aria-hidden />

      <div className="glass-toolbar mac-no-drag flex min-w-0 flex-1 items-center gap-2 py-1.5 pr-1.5 pl-4">
        <IconSearch
          className="size-4 shrink-0 text-muted-foreground"
          stroke={1.7}
        />
        <input
          ref={searchRef}
          type="search"
          role="searchbox"
          aria-label="Search devices, sessions, and settings"
          placeholder="Search"
          onChange={(e) => onSearch?.(e.target.value)}
          className="mac-focus min-w-0 flex-1 bg-transparent text-[13px] outline-none placeholder:text-muted-foreground"
        />
        <Kbd className="hidden shrink-0 sm:inline-flex">⌘K</Kbd>

        <span
          role="status"
          className={cn(
            "hidden shrink-0 items-center gap-1.5 rounded-full px-2.5 py-1 text-[11px] font-medium md:inline-flex",
            online
              ? "bg-primary/15 text-primary"
              : "bg-foreground/[0.06] text-muted-foreground",
          )}
        >
          <IconCircleFilled
            className={cn(
              "size-1.5",
              online ? "text-primary" : "text-muted-foreground/50",
            )}
          />
          {online ? "Online" : "Offline"}
        </span>

        <Button
          variant="ghost"
          size="icon-sm"
          onClick={toggleTheme}
          title={`Theme: ${mode} (click to toggle)`}
          aria-label={`Theme: ${mode}. Activate to toggle.`}
          className="mac-press shrink-0 rounded-full"
        >
          {resolved === "dark" ? (
            <IconSun className="size-4" stroke={1.7} />
          ) : (
            <IconMoon className="size-4" stroke={1.7} />
          )}
        </Button>

        <Button
          onClick={onConnect}
          className="mac-capsule mac-press glass-prominent h-8 shrink-0 border px-4 text-[13px] font-semibold text-foreground hover:opacity-90"
        >
          Connect
        </Button>
      </div>
    </div>
  );
}
