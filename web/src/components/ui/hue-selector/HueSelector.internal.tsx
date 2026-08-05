import { CSSProperties } from "react";

import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import { styled } from "@/styled-system/jsx";

import { C, type HueSelectorProps, L, useHueSelector } from "./useHueSelector";

export function HueSelector(props: HueSelectorProps) {
  const { onPointerDown, onPointerUp, hue, ref, angle, value, grabbing } =
    useHueSelector(props);

  const styles = {
    backgroundColor: value,
    "--angle": `${angle}deg`,
    "--thumb-size": "3em",
    "--circle-size": "200px",
    "--colour": value,
    "--cursor": grabbing ? "grabbing" : "grab",
  } as CSSProperties;

  return (
    <styled.div
      width="min"
      style={styles}
      borderWidth="thin"
      borderStyle="solid"
      borderColor="border.strong"
      borderRadius="lg"
      p="4"
    >
      <styled.div
        ref={ref}
        borderRadius="full"
        position="relative"
        background="conicGradient"
        style={{
          touchAction: "none",
          width: "var(--circle-size)",
          height: "var(--circle-size)",
        }}
        boxShadow="surface"
        _before={
          {
            content: '""',
            position: "absolute",
            top: "var(--thumb-size)",
            left: "var(--thumb-size)",
            alignItems: "center",
            justifyContent: "center",
            borderRadius: "50%",
            display: "flex",
            height: "calc(var(--circle-size) - (var(--thumb-size) * 2))",
            width: "calc(var(--circle-size) - (var(--thumb-size) * 2))",
            backgroundColor: "var(--colour)",
          } as any
        }
      >
        <styled.output
          position="absolute"
          display="flex"
          justifyContent="end"
          alignItems="center"
          transform="rotate(var(--angle))"
          transformOrigin="center left"
          onPointerDown={onPointerDown}
          onPointerUp={onPointerUp}
          cursor="var(--cursor)"
          style={{
            width: "50%",
            height: "var(--thumb-size)",
            top: "50%",
            left: "50%",
            marginTop: "calc(var(--thumb-size) / -2)",
          }}
        >
          <styled.div
            display="flex"
            justifyContent="center"
            alignItems="center"
            borderRadius="lg"
            backgroundColor="white"
            height="8"
            style={{
              width: "var(--thumb-size)",
            }}
          >
            <DragHandleIcon color="text.subtle" />
          </styled.div>
        </styled.output>
      </styled.div>
    </styled.div>
  );
}
