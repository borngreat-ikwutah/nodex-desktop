/**
 * OLED substrate for the console: two slow mesh orbs and a film grain veil.
 * Fixed and pointer-events-none so scrolling content never repaints the blur.
 */
export function MeshBackdrop() {
  return (
    <div
      aria-hidden
      className="pointer-events-none fixed inset-0 z-0 overflow-hidden"
    >
      <div className="absolute -top-40 -left-32 size-[42rem] rounded-full bg-[oklch(0.72_0.16_163/0.16)] blur-[120px]" />
      <div className="absolute -right-40 top-1/3 size-[38rem] rounded-full bg-[oklch(0.62_0.19_292/0.14)] blur-[130px]" />
      <div className="absolute bottom-0 left-1/3 size-[30rem] rounded-full bg-[oklch(0.75_0.14_200/0.08)] blur-[140px]" />
      <div className="absolute inset-0 opacity-[0.035] mix-blend-overlay [background-image:url('data:image/svg+xml;utf8,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22160%22 height=%22160%22><filter id=%22n%22><feTurbulence type=%22fractalNoise%22 baseFrequency=%220.85%22 numOctaves=%223%22/></filter><rect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23n)%22/></svg>')]" />
    </div>
  );
}
