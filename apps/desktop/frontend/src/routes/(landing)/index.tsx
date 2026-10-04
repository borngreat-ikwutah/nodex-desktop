import { createFileRoute } from "@tanstack/react-router";
import { LandingPage } from "../../features/landing/landing-page";

export const Route = createFileRoute("/(landing)/")({
  component: RouteComponent,
});

function RouteComponent() {
  return <LandingPage />;
}
