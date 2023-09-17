import Color from "colorjs.io";
import { range } from "lodash";
import { map } from "lodash/fp";
import { useEffect, useRef, useState } from "react";

export type Props = {
  onChange: (value: string) => void;
  onUpdate: (value: string) => void;
  value: string;
};

// TODO: Dark mode = 40%
export const L = "80%";

export const C = "0.15";

export const lch = (hue: number) => `oklch(${L} ${C} ${hue})`;

const stops = map(lch)(range(0, 361, 10));

export const conicGradient = `
conic-gradient(
    ${stops.join(",\n")}
);
`;

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
      const colour = new Color(props.value);

      const hue = angleToHue(colour.lch["h"] ?? 0);

      if (hue) {
        setAngle(hue);
      }
    } catch (_) {
      setAngle(Math.random() * 359);
    }
  }, [props.value]);

  const hue = hueToAngle(angle);
  const value = `oklch(${L} ${C} ${hue}deg)`;

  function onPointerMove(e: globalThis.PointerEvent) {
    if (!ref.current) return;

    const rect = ref.current.getBoundingClientRect();

    const cx = rect.left + rect.width / 2;
    const cy = rect.top + rect.height / 2;

    const mx = e.clientX;
    const my = e.clientY;

    const angleTo = (Math.atan2(my - cy, mx - cx) * 180) / Math.PI;

    setAngle(angleTo);
    props.onUpdate(`oklch(${L} ${C} ${hueToAngle(angleTo)}deg)`);
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
