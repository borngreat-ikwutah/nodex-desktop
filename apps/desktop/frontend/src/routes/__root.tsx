import { createRootRoute } from "@tanstack/react-router";
import "../app.css";
import { AppLayout } from "../components/app-layout";

export const Route = createRootRoute({
  component: AppLayout,
  errorComponent: () => (
    <div className="p-4 text-destructive">
      Something went wrong at the root!
    </div>
  ),
  pendingComponent: () => <div className="p-4">Loading root...</div>,
});
