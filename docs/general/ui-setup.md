# UI Setup & Requirements

## Shadcn Components

### Tooltip Provider

The `tooltip` component requires your application root to be wrapped with the `TooltipProvider`.

When integrating `@repo/ui` into your apps, make sure to add this provider to your root layout (e.g., `App.tsx` or `layout.tsx`):

```tsx
import { TooltipProvider } from "@repo/ui/components/ui/tooltip";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <TooltipProvider>{children}</TooltipProvider>;
}
```
