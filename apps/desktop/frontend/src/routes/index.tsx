import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: Index,
  loader: async () => {
    await new Promise((r) => setTimeout(r, 200));
    return { message: "Welcome to nodex-desktop!" };
  },
});

function Index() {
  const { message } = Route.useLoaderData();

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground mt-2">{message}</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {/* Placeholder cards to show theme adherence */}
        {Array.from({ length: 4 }).map((_, i) => (
          <div
            key={i}
            className="rounded-xl border bg-card text-card-foreground shadow-sm"
          >
            <div className="p-6 flex flex-row items-center justify-between space-y-0 pb-2">
              <h3 className="tracking-tight text-sm font-medium">
                Stat {i + 1}
              </h3>
            </div>
            <div className="p-6 pt-0">
              <div className="text-2xl font-bold">+123</div>
              <p className="text-xs text-muted-foreground">
                +19% from last month
              </p>
            </div>
          </div>
        ))}
      </div>

      <div className="rounded-xl border bg-card text-card-foreground shadow-sm">
        <div className="p-6">
          <h3 className="text-lg font-semibold leading-none tracking-tight">
            Recent Activity
          </h3>
          <p className="text-sm text-muted-foreground mt-2">
            Your recent application activity.
          </p>
        </div>
        <div className="p-6 pt-0">
          <div className="space-y-8">
            <div className="flex items-center">
              <div className="ml-4 space-y-1">
                <p className="text-sm font-medium leading-none">
                  System Started
                </p>
                <p className="text-sm text-muted-foreground">
                  The desktop engine is running smoothly.
                </p>
              </div>
              <div className="ml-auto font-medium text-sm text-muted-foreground">
                Just now
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
