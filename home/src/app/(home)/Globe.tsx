"use client";

import { Box, BoxProps, Center } from "@/styled-system/jsx";
import createGlobe from "cobe";
import { useRef, useEffect, useState } from "react";

const SIZE = 250;

export function Globe(props: BoxProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const sizeRef = useRef(SIZE);

  function handleResize() {
    if (!containerRef.current) {
      return;
    }

    const { clientWidth, clientHeight } = containerRef.current;

    sizeRef.current = Math.min(clientWidth, clientHeight);
  }

  useEffect(() => {
    let phi = 0;

    if (!canvasRef.current || !containerRef.current) {
      return;
    }

    handleResize();
    const size = sizeRef.current;

    const globe = createGlobe(canvasRef.current, {
      devicePixelRatio: 2,
      width: size * 2,
      height: size * 2,
      phi: 0,
      theta: 0.18,
      dark: 0,
      diffuse: 0,
      mapSamples: 6000,
      mapBrightness: 10,
      mapBaseBrightness: 0,
      baseColor: [1, 1, 1],
      markerColor: [1, 1, 1],
      glowColor: [1, 1, 1],
      opacity: 0.92,
      markers: [],
      onRender: (state) => {
        // Called on every animation frame.
        // `state` will be an empty object, return updated params.
        state.phi = phi;
        phi += 0.001;
        state.width = sizeRef.current * 2;
        state.height = sizeRef.current * 2;
      },
    });

    window.addEventListener("resize", handleResize);

    return () => {
      globe.destroy();
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  return (
    <Center ref={containerRef} {...props}>
      <canvas
        ref={canvasRef}
        style={{
          width: SIZE,
          height: SIZE,
          maxWidth: "100%",
          aspectRatio: 1,
        }}
      />
    </Center>
  );
}
