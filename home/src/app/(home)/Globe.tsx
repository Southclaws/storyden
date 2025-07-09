"use client";

import { BoxProps, Center } from "@/styled-system/jsx";
import createGlobe from "cobe";
import { useEffect, useRef, useState } from "react";

// An initial size for most laptop screens. Flashes but is below the fold so meh
const INITIAL_SIZE = 125;

export function Globe(props: BoxProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  // Used for constraining the <canvas> to the exact size of its parent as well
  // as locking its aspect ratio to the shortest edge. Not used for the actual
  // canvas API sizing, just the DOM context sizing.
  const [containerSize, setContainerSize] = useState<number>(INITIAL_SIZE);

  function handleResize() {
    if (!canvasRef.current || !containerRef.current) {
      return;
    }

    const dprScale = window.devicePixelRatio ?? 1;

    const { clientWidth, clientHeight } = containerRef.current;

    const constrainedSize = Math.min(clientWidth, clientHeight);

    setContainerSize(constrainedSize);

    // Use the device pixel ratio to make it crispy.
    canvasRef.current.width = constrainedSize * dprScale;
    canvasRef.current.height = constrainedSize * dprScale;
  }

  let phi = 0;

  useEffect(() => {
    if (!canvasRef.current || !containerRef.current) {
      return;
    }

    // Calculate initial size for canvasRef values.
    handleResize();

    const globe = createGlobe(canvasRef.current, {
      devicePixelRatio: 2,
      width: canvasRef.current.width,
      height: canvasRef.current.height,
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
        if (!canvasRef.current) {
          return;
        }

        // Rotate slowly.
        state.phi = phi;
        phi += 0.001;

        // Update width/height to reflect actual canvas size.
        state.width = canvasRef.current.width;
        state.height = canvasRef.current.height;
      },
    });

    const resizeObserver = new ResizeObserver(() => {
      handleResize();
      globe.resize();
    });

    resizeObserver.observe(containerRef.current);

    return () => {
      globe.destroy();
      resizeObserver.disconnect();
    };
  }, []);

  return (
    <Center ref={containerRef} w="full" h="full" {...props}>
      <canvas
        ref={canvasRef}
        style={{
          width: containerSize,
          height: containerSize,
          maxWidth: "100%",
          aspectRatio: 1,
        }}
      />
    </Center>
  );
}
