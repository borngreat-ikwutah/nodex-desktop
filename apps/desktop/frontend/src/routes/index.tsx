import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: Index,
  loader: async () => {
    // Intent preloading and route loader example
    await new Promise(r => setTimeout(r, 200))
    return { message: 'Welcome to nodex-desktop!' }
  }
})

import { Button } from '@repo/ui/components/ui/button'

function Index() {
  const { message } = Route.useLoaderData()
  return (
    <div className="p-4 flex flex-col items-start gap-4">
      <h3 className="text-xl font-bold">{message}</h3>
      <Button onClick={() => alert('Clicked!')}>Click Me</Button>
    </div>
  )
}
