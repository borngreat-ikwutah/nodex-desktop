import { Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { Suspense } from "react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarInset,
} from "@repo/ui/components/ui/sidebar.tsx";
import { TooltipProvider } from "@repo/ui/components/ui/tooltip.tsx";
import { IconDiamondFilled, IconCircle } from "@tabler/icons-react";
import { navigationData } from "../data/navigation";

function TitleBar() {
  return (
    <div className="flex h-12 shrink-0 items-center justify-between border-b bg-background px-4">
      <div className="flex items-center gap-2 font-semibold">
        <IconDiamondFilled className="size-5 text-primary" />
        <span>Relay</span>
      </div>
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <IconCircle className="size-3" />
        <span>Offline</span>
      </div>
    </div>
  );
}

export function AppLayout() {
  return (
    <TooltipProvider>
      <div className="flex h-screen flex-col overflow-hidden bg-background text-foreground">
        <TitleBar />
        <SidebarProvider className="flex-1 overflow-hidden">
          <Sidebar className="border-r">
            <SidebarContent className="py-4">
              <SidebarGroup>
                <SidebarGroupContent>
                  <SidebarMenu>
                    {navigationData.map((item) => (
                      <SidebarMenuItem key={item.title}>
                        <SidebarMenuButton asChild>
                          <Link
                            to={item.url as any}
                            className="[&.active]:bg-sidebar-accent [&.active]:text-sidebar-accent-foreground"
                          >
                            <item.icon className="size-4" />
                            <span>{item.title}</span>
                          </Link>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    ))}
                  </SidebarMenu>
                </SidebarGroupContent>
              </SidebarGroup>
            </SidebarContent>
          </Sidebar>
          <SidebarInset className="flex-1 overflow-hidden bg-transparent">
            {/* The tailwind demo grid layout */}
            <div className="relative grid h-full grid-cols-[1fr_2.5rem_auto_2.5rem_1fr] grid-rows-[1fr_1px_auto_1px_1fr] bg-white [--pattern-fg:var(--color-gray-950)]/5 dark:bg-gray-950 dark:[--pattern-fg:var(--color-white)]/10">
              <div className="col-start-3 row-start-3 flex w-full min-w-[50vw] max-w-5xl flex-col bg-gray-100 p-2 dark:bg-white/10">
                <div className="flex-1 rounded-xl bg-white p-8 text-sm/7 text-gray-700 shadow-sm dark:bg-gray-950 dark:text-gray-300">
                  <Suspense
                    fallback={<div className="p-4">Loading Workspace...</div>}
                  >
                    <Outlet />
                  </Suspense>
                </div>
              </div>
              <div className="relative -right-px col-start-2 row-span-full row-start-1 border-x border-x-(--pattern-fg) bg-[image:repeating-linear-gradient(315deg,_var(--pattern-fg)_0,_var(--pattern-fg)_1px,_transparent_0,_transparent_50%)] bg-[size:10px_10px] bg-fixed"></div>
              <div className="relative -left-px col-start-4 row-span-full row-start-1 border-x border-x-(--pattern-fg) bg-[image:repeating-linear-gradient(315deg,_var(--pattern-fg)_0,_var(--pattern-fg)_1px,_transparent_0,_transparent_50%)] bg-[size:10px_10px] bg-fixed"></div>
              <div className="relative -bottom-px col-span-full col-start-1 row-start-2 h-px bg-(--pattern-fg)"></div>
              <div className="relative -top-px col-span-full col-start-1 row-start-4 h-px bg-(--pattern-fg)"></div>
            </div>
          </SidebarInset>
        </SidebarProvider>
      </div>
      <TanStackRouterDevtools />
    </TooltipProvider>
  );
}
