import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/transfers")({
  component: () => (
    <div className="flex h-full items-center justify-center text-muted-foreground">
      transfers workspace (empty for now)
    </div>
  ),
});
