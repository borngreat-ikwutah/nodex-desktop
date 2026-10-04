import { Button } from "@repo/ui/components/ui/button.tsx";
import { IconCopy, IconEye, IconDeviceDesktop } from "@tabler/icons-react";
import logoUrl from "../../assets/logo.svg";

export function HomePage() {
  return (
    <div className="flex-1 rounded-xl bg-card p-6 lg:p-8 text-card-foreground shadow-sm border overflow-auto flex flex-col w-full h-full">
      <div className="flex items-center gap-3 border-b pb-4 mb-6">
        <img
          src={logoUrl}
          alt="Nodex Logo"
          className="w-6 h-auto drop-shadow-sm text-primary"
        />
        <h1 className="text-lg font-bold tracking-tight text-foreground">
          Workspace Overview
        </h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 flex-1">
        {/* This Device */}
        <div className="flex flex-col gap-3">
          <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Allow Remote Control
          </h2>
          <div className="rounded-lg border bg-background p-5 shadow-sm space-y-5">
            <div className="flex items-center justify-between pb-4 border-b border-border/50">
              <div className="flex items-center gap-2 text-sm font-medium">
                <IconDeviceDesktop className="size-4 text-muted-foreground" />
                <span>This Device</span>
              </div>
              <span className="flex items-center gap-2 text-xs font-medium text-primary">
                <span className="relative flex h-2 w-2">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
                </span>
                Ready to Connect
              </span>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-medium text-muted-foreground">
                Your Nodex ID
              </label>
              <div className="text-xl font-mono font-semibold tracking-wider bg-muted/50 p-2.5 rounded-md flex justify-between items-center border">
                <span>123 456 789</span>
                <Button
                  variant="ghost"
                  size="icon-sm"
                  className="h-7 w-7 text-muted-foreground hover:text-foreground"
                >
                  <IconCopy className="size-4" />
                </Button>
              </div>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-medium text-muted-foreground">
                Session Password
              </label>
              <div className="text-lg font-mono tracking-widest bg-muted/50 p-2.5 rounded-md flex justify-between items-center border">
                <span>••••••••</span>
                <div className="flex items-center gap-1">
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    className="h-7 w-7 text-muted-foreground hover:text-foreground"
                  >
                    <IconEye className="size-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    className="h-7 w-7 text-muted-foreground hover:text-foreground"
                  >
                    <IconCopy className="size-4" />
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Remote Session */}
        <div className="flex flex-col gap-3">
          <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Control Remote Device
          </h2>
          <div className="rounded-lg border bg-background p-5 shadow-sm space-y-4">
            <p className="text-sm text-muted-foreground leading-relaxed">
              Enter a partner's Nodex ID to view or control their device.
            </p>
            <div className="space-y-3 pt-1">
              <input
                type="text"
                placeholder="Partner ID"
                className="flex h-11 w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 font-mono tracking-wider"
              />
              <Button className="w-full h-11 font-semibold" size="lg">
                Connect
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
