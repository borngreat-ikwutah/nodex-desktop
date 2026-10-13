import type { ReactNode } from "react";
import { ThemeProvider } from "../lib/theme";

/** Native vibrancy root + theme scope. Replaces the opaque bg-background shell. */
export function MacWindow({ children }: { children: ReactNode }) {
  return (
    <ThemeProvider>
      <div className="mac-vibrancy flex h-screen flex-col overflow-hidden text-foreground select-none">
        {children}
      </div>
    </ThemeProvider>
  );
}
