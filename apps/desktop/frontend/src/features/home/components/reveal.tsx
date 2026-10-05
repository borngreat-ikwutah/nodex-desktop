import { useEffect, useRef, useState } from "react";
import type { ReactNode } from "react";

type RevealProps = {
  children: ReactNode;
  delay?: number;
  threshold?: number;
  className?: string;
};

const HIDDEN = "opacity-0 translate-y-6 blur-md";
const VISIBLE = "opacity-100 translate-y-0 blur-0";

/**
 * Reveals its children once they scroll into view. The transition animates only
 * transform, opacity and filter, so entry never triggers layout in the grid.
 */
export function Reveal({
  children,
  delay = 0,
  threshold = 0.15,
  className = "",
}: RevealProps) {
  const ref = useRef<HTMLDivElement | null>(null);
  const [isVisible, setIsVisible] = useState(
    () => typeof IntersectionObserver === "undefined",
  );

  useEffect(() => {
    const node = ref.current;
    if (!node || isVisible) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (!entries.some((entry) => entry.isIntersecting)) return;
        setIsVisible(true);
        observer.disconnect();
      },
      { threshold },
    );

    observer.observe(node);
    return () => observer.disconnect();
  }, [isVisible, threshold]);

  return (
    <div
      ref={ref}
      style={{
        transitionDelay: `${delay}ms`,
        transitionProperty: "opacity, transform, filter",
        transitionTimingFunction: "cubic-bezier(0.32, 0.72, 0, 1)",
        transitionDuration: "900ms",
      }}
      className={`${isVisible ? VISIBLE : HIDDEN} ${className}`}
    >
      {children}
    </div>
  );
}
