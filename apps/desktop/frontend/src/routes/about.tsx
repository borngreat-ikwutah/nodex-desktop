import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/about')({
  component: About,
  validateSearch: (search: Record<string, unknown>): { debug?: boolean } => {
    // Typed search schemas example
    return {
      debug: Boolean(search?.debug || false),
    }
  },
})

function About() {
  const { debug } = Route.useSearch()
  return (
    <div className="p-2">
      <p>About nodex-desktop</p>
      {debug && <p>Debug mode is ON</p>}
    </div>
  )
}
