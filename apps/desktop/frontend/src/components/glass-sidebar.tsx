import { Link, useRouterState } from "@tanstack/react-router";
import { IconDiamondFilled } from "@tabler/icons-react";
import { navigationData } from "../data/navigation";
import { cn } from "../lib/cn";

/**
 * Inset floating glass sidebar (NSSplitView sidebar analogue).
 * Expanded ≥1024px, icon rail on smaller widths via CSS.
 */
export function GlassSidebar() {
  const pathname = useRouterState({ select: (s) => s.location.pathname });

  return (
    <nav
      aria-label="Primary"
      className="glass-sidebar mac-no-drag hidden w-[68px] shrink-0 flex-col gap-1 overflow-y-auto p-2 sm:flex lg:w-[232px] lg:p-3"
    >
      <div className="mb-2 flex items-center gap-2.5 px-2 pt-1.5 pb-1">
        <span className="glass-prominent flex size-8 shrink-0 items-center justify-center rounded-[10px]">
          <IconDiamondFilled className="size-4 text-primary" stroke={1.7} />
        </span>
        <span className="hidden text-[13px] font-semibold tracking-tight lg:block">
          Nodex
        </span>
      </div>

      {navigationData.map((item) => {
        const active =
          pathname === item.url || pathname.startsWith(`${item.url}/`);
        return (
          <Link
            key={item.title}
            to={item.url as never}
            title={item.title}
            aria-current={active ? "page" : undefined}
            className={cn(
              "mac-press mac-focus group flex items-center justify-center gap-3 rounded-[10px] px-2 py-2 text-[13px] font-medium lg:justify-start lg:px-3",
              active
                ? "glass-prominent text-foreground"
                : "text-muted-foreground hover:bg-foreground/[0.06] hover:text-foreground",
            )}
          >
            <item.icon
              className={cn(
                "size-[17px] shrink-0",
                active ? "text-primary" : "text-muted-foreground",
              )}
              stroke={1.7}
            />
            <span className="hidden lg:block">{item.title}</span>
            {active && (
              <span
                aria-hidden
                className="ml-auto hidden size-1.5 rounded-full bg-primary lg:block"
              />
            )}
          </Link>
        );
      })}

      <div className="mt-auto hidden px-2 pt-3 pb-1 lg:block">
        <p className="text-[11px] leading-relaxed text-muted-foreground">
          Liquid Glass · follows system theme
        </p>
      </div>
    </nav>
  );
}
