import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/sessions")({
  component: () => (
    <div className="flex h-full items-center justify-center text-muted-foreground">
      sessions workspace (empty for now)
    </div>
  ),
});
