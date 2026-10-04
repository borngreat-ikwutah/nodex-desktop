import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/devices")({
  component: () => (
    <div className="flex h-full items-center justify-center text-muted-foreground">
      devices workspace (empty for now)
    </div>
  ),
});
