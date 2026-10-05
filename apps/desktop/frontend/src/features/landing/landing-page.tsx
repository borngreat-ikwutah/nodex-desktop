import { Link } from "@tanstack/react-router";
import { Button } from "@repo/ui/components/ui/button.tsx";
import logoUrl from "../../assets/logo.svg";

export function LandingPage() {
  return (
    <div className="flex-1 rounded-xl bg-card p-4 lg:p-8 text-card-foreground shadow-sm border overflow-auto flex flex-col items-center justify-center relative w-full h-full">
      <div className="z-10 flex flex-col items-center text-center max-w-2xl space-y-4 lg:space-y-8 mt-0 lg:mt-[-10vh]">
        <div className="flex items-center gap-3 lg:gap-4">
          <img
            src={logoUrl}
            alt="Nodex Logo"
            className="w-10 lg:w-16 h-auto drop-shadow-md text-primary"
          />
          <h1 className="text-3xl lg:text-6xl font-extrabold tracking-tight text-foreground drop-shadow-sm">
            Nodex
          </h1>
        </div>

        <p className="text-base lg:text-xl text-muted-foreground leading-relaxed">
          The ultimate remote desktop experience. Connect, manage, and
          collaborate securely across all your devices in real-time.
        </p>

        <div className="flex flex-col lg:flex-row items-center gap-3 lg:gap-4 pt-2 lg:pt-4 w-full lg:w-auto">
          <Button
            asChild
            size="lg"
            className="h-10 lg:h-12 px-6 lg:px-8 text-sm lg:text-base font-semibold shadow-lg shadow-primary/20 w-full lg:w-auto"
          >
            <Link to="/home">Get Started</Link>
          </Button>
          <Button
            size="lg"
            variant="outline"
            className="h-10 lg:h-12 px-6 lg:px-8 text-sm lg:text-base font-semibold w-full lg:w-auto"
          >
            Read the Docs
          </Button>
        </div>
      </div>
    </div>
  );
}
