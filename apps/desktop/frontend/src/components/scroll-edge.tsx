/** Hard-edge scroll separation strip used under floating glass bars (macOS style). */
export function ScrollEdge({ className = "" }: { className?: string }) {
  return (
    <div
      aria-hidden
      className={`mac-scroll-fade pointer-events-none sticky top-0 z-10 h-7 ${className}`}
    />
  );
}
