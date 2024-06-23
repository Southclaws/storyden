"use client";

import { parseToHsl } from "polished";
import { useEffect, useRef, useState } from "react";

export type Props = {
  onChange: (value: string) => void;
  onUpdate: (value: string) => void;
  value: string;
};

// TODO: Dark mode = 40%
export const L = "80";

export const C = "15";

// TODO: Support LCH where supported.
// export const lch = (hue: number) => `oklch(${L} ${C} ${hue})`;
export const hsl = (hue: number) => `hsl(${hue}, ${C}%, ${L}%)`;

const hueToAngle = (input: number) => {
  const shifted = input + 90;

  let clamped = shifted;
  while (clamped < 0) clamped += 360;

  return clamped;
};

const angleToHue = (input: number) => {
  const shifted = input - 90;

  let clamped = shifted;
  while (clamped < 0) clamped += 360;

  return clamped;
};

export function useColourInput(props: Props) {
  const ref = useRef<HTMLDivElement>(null);
  const [angle, setAngle] = useState(270);
  const [grabbing, setGrabbing] = useState(false);

  useEffect(() => {
    if (!props.value) return;

    try {
      const colour = parseToHsl(props.value);

      const hue = angleToHue(colour.hue ?? 0);

      if (hue) {
        setAngle(hue);
      }
    } catch (e) {
      setAngle(Math.random() * 359);
    }
  }, [props.value]);

  const hue = hueToAngle(angle);
  const value = hsl(hue);

  function onPointerMove(e: globalThis.PointerEvent) {
    if (!ref.current) return;

    const rect = ref.current.getBoundingClientRect();

    const cx = rect.left + rect.width / 2;
    const cy = rect.top + rect.height / 2;

    const mx = e.clientX;
    const my = e.clientY;

    const angleTo = (Math.atan2(my - cy, mx - cx) * 180) / Math.PI;

    setAngle(angleTo);
    props.onUpdate(hsl(hueToAngle(angleTo)));
  }

  function onCleanup() {
    setGrabbing(false);
    removeEventListener("pointermove", onPointerMove);
    removeEventListener("pointerup", onCleanup);
  }

  function onPointerDown() {
    setGrabbing(true);
    addEventListener("pointermove", onPointerMove);
    addEventListener("pointerup", onCleanup);
  }

  function onPointerUp() {
    onCleanup();
    props.onChange(value);
  }

  return {
    onPointerDown,
    onPointerUp,
    hue,
    ref,
    angle,
    value,
    grabbing,
  };
}
