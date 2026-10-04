import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/(landing)/")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex h-full flex-col">
      <h1 className="text-2xl font-bold tracking-tight">Home</h1>
      <div className="flex flex-1 items-center justify-center text-muted-foreground">
        Main Workspace (empty for now - just placeholders)
      </div>
    </div>
  );
}
