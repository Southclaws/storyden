"use client";

import { motion, useReducedMotion } from "framer-motion";
import { useState } from "react";

import { styled } from "@/styled-system/jsx";

export type PinchGestureHintProps = {
  size?: number;
  className?: string;
  onComplete?: () => void;
};

const times = [0, 0.11, 0.47, 0.525, 0.635, 1];
const offsets = [8, 8, 20, 20, 20, 8];
const fadeEase = [0.42, 0, 1, 1] as const;

/** A one-shot gesture cue, or a static illustration for reduced motion. */
export function PinchGestureHint({
  size = 72,
  className,
  onComplete,
}: PinchGestureHintProps) {
  const reducedMotion = useReducedMotion();
  const [complete, setComplete] = useState(false);
  const animate = reducedMotion === false;

  return (
    <styled.span
      className={className}
      display="inline-flex"
      flexShrink={0}
      pointerEvents="none"
    >
      <styled.span srOnly>Use two fingers to zoom</styled.span>
      {(!complete || !animate) && (
        <motion.svg
          aria-hidden="true"
          focusable="false"
          width={size}
          height={size}
          viewBox="0 0 64 64"
          fill="currentColor"
          fillOpacity={0.28}
          stroke="currentColor"
          strokeWidth={1.5}
          style={{ pointerEvents: "none", flexShrink: 0 }}
          initial={animate ? { opacity: 0 } : false}
          animate={animate ? { opacity: [0, 1, 0, 0, 1, 0] } : undefined}
          transition={{
            duration: 2,
            times,
            ease: ["linear", fadeEase, "linear", "linear", fadeEase],
          }}
          onAnimationComplete={() => {
            setComplete(true);
            onComplete?.();
          }}
        >
          {!animate && (
            <path
              d="M19 45 45 19 M19 37v8h8 M37 19h8v8"
              fill="none"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          )}
          {[-1, 1].map((direction) => (
            <motion.circle
              key={direction}
              cx={32 + direction * 20}
              cy={32 - direction * 20}
              r={5}
              initial={
                animate ? { x: direction * -12, y: direction * 12 } : false
              }
              animate={
                animate
                  ? {
                      x: offsets.map((offset) => direction * (offset - 20)),
                      y: offsets.map((offset) => direction * (20 - offset)),
                    }
                  : undefined
              }
              transition={{
                duration: 2,
                times,
                ease: [0.22, 1, 0.36, 1],
              }}
            />
          ))}
        </motion.svg>
      )}
    </styled.span>
  );
}
