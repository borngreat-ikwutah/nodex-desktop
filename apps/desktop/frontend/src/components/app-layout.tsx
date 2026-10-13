import { Outlet, useNavigate } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { Suspense, useCallback } from "react";
import { TooltipProvider } from "@repo/ui/components/ui/tooltip.tsx";
import { Toaster } from "sonner";
import { MacWindow } from "./mac-window";
import { GlassToolbar } from "./glass-toolbar";
import { GlassSidebar } from "./glass-sidebar";
import { useTheme } from "../lib/theme";

/** Floating glass app shell: native toolbar + inset sidebar + edge-to-edge content. */
export function AppLayout() {
  const navigate = useNavigate();
  const { resolved } = useTheme();

  const handleConnect = useCallback(() => {
    navigate({ to: "/home" as never });
    // Home owns the connect sheet; toolbar deep-links via hash event.
    window.dispatchEvent(new CustomEvent("nodex:open-connect"));
  }, [navigate]);

  return (
    <TooltipProvider>
      <MacWindow>
        <GlassToolbar status="offline" onConnect={handleConnect} />
        <div className="flex min-h-0 flex-1 gap-3 px-3 pt-3 pb-3 sm:px-4 sm:pb-4">
          <GlassSidebar />
          <main
            aria-label="Workspace"
            className="relative min-h-0 min-w-0 flex-1 overflow-y-auto rounded-[var(--mac-window-radius)]"
          >
            <Suspense
              fallback={
                <div className="flex h-full items-center justify-center p-4 text-[13px] text-muted-foreground">
                  Loading workspace…
                </div>
              }
            >
              <Outlet />
            </Suspense>
          </main>
        </div>
      </MacWindow>
      <Toaster theme={resolved === "dark" ? "dark" : "light"} />
      <TanStackRouterDevtools />
    </TooltipProvider>
  );
}
