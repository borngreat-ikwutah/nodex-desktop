import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { Suspense } from "react";
import "../app.css";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
  SidebarInset,
  SidebarHeader,
} from "@repo/ui/components/ui/sidebar.tsx";

import { TooltipProvider } from "@repo/ui/components/ui/tooltip.tsx";

import { IconHome, IconInfoCircle } from "@tabler/icons-react";

export const Route = createRootRoute({
  component: AppSidebar,
  errorComponent: () => (
    <div className="p-4 text-destructive">
      Something went wrong at the root!
    </div>
  ),
  pendingComponent: () => <div className="p-4">Loading root...</div>,
});

function AppSidebar() {
  return (
    <TooltipProvider>
      <SidebarProvider>
        <Sidebar>
          <SidebarHeader>
            <div className="p-4 font-bold text-lg">Nodex Desktop</div>
          </SidebarHeader>
          <SidebarContent>
            <SidebarGroup>
              <SidebarGroupLabel>Navigation</SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link
                        to="/"
                        className="[&.active]:bg-sidebar-accent [&.active]:text-sidebar-accent-foreground"
                      >
                        <IconHome className="size-4" />
                        <span>Home</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                  <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                      <Link
                        to="/about"
                        className="[&.active]:bg-sidebar-accent [&.active]:text-sidebar-accent-foreground"
                      >
                        <IconInfoCircle className="size-4" />
                        <span>About</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </SidebarContent>
        </Sidebar>
        <SidebarInset>
          <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
            <SidebarTrigger />
          </header>
          <main className="flex flex-1 flex-col p-4 bg-background text-foreground">
            <Suspense fallback={<div className="p-4">Loading...</div>}>
              <Outlet />
            </Suspense>
          </main>
        </SidebarInset>
        <TanStackRouterDevtools />
      </SidebarProvider>
    </TooltipProvider>
  );
}
