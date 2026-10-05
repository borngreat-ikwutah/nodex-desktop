import { createRootRoute } from "@tanstack/react-router";
import "../app.css";
import { AppLayout } from "../components/app-layout";

export const Route = createRootRoute({
  component: AppLayout,
  errorComponent: ({ error }) => (
    <div className="p-4 text-destructive">
      Something went wrong at the root! {(error as Error)?.message}
    </div>
  ),
  pendingComponent: () => <div className="p-4">Loading root...</div>,
});
