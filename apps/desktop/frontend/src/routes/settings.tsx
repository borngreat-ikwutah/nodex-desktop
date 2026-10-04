import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/settings")({
  component: () => (
    <div className="flex h-full items-center justify-center text-muted-foreground">
      settings workspace (empty for now)
    </div>
  ),
});
