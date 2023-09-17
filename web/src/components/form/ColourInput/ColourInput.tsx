import { Box, FormControl, forwardRef } from "@chakra-ui/react";
import { CSSProperties } from "react";
import { Control, Controller } from "react-hook-form";

import { DragHandleIcon } from "../../graphics/DragHandleIcon";

import { C, L, Props, conicGradient, useColourInput } from "./useColourInput";

export function ColourInput(props: Props) {
  const { onPointerDown, onPointerUp, hue, ref, angle, value, grabbing } =
    useColourInput(props);

  // NOTE: we use React's primitive `style` prop here and not "style props" from
  // Chakra because this is going to change on every frame while dragging and we
  // don't want Emotion.js to generate a whole new stylesheet when these change!
  const styles = {
    backgroundColor: value,
    "--angle": `${angle}deg`,
    "--thumb-size": "3em",
    "--circle-size": "200px",
    "--colour": `oklch(${L} ${C} ${hue}deg)`,
    "--cursor": grabbing ? "grabbing" : "grab",
  } as CSSProperties;

  return (
    <Box
      width="min-content"
      style={styles}
      borderWidth={1}
      borderStyle="solid"
      borderColor="blackAlpha.100"
      borderRadius={10}
      p={4}
    >
      <Box
        ref={ref}
        width="var(--circle-size)"
        height="var(--circle-size)"
        borderRadius="full"
        position="relative"
        background={conicGradient}
        style={{
          touchAction: "none",
        }}
        boxShadow="0 10px 30px rgba(0, 0, 0, 0.05)"
        _before={{
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
        }}
      >
        <Box
          as="output"
          position="absolute"
          width="50%"
          height="var(--thumb-size)"
          top="50%"
          left="50%"
          display="flex"
          justifyContent="end"
          alignItems="center"
          transform="rotate(var(--angle))"
          transformOrigin="center left"
          marginTop="calc(var(--thumb-size) / -2)"
          onPointerDown={onPointerDown}
          onPointerUp={onPointerUp}
          cursor="var(--cursor)"
        >
          <Box
            display="flex"
            justifyContent="center"
            alignItems="center"
            borderRadius="lg"
            backgroundColor="white"
            width="var(--thumb-size)"
            height="2em"
          >
            <DragHandleIcon />
          </Box>
        </Box>
      </Box>
    </Box>
  );
}

type FieldProps = {
  control: Control<any>;
  name: string;
  defaultValue: string;
  onUpdate: (v: string) => void;
};

export const ColourField = forwardRef((props: FieldProps, ref) => {
  return (
    <FormControl ref={ref}>
      <Controller
        defaultValue={props.defaultValue}
        render={({ field: { onChange, ...field } }) => {
          return (
            <ColourInput
              onChange={onChange}
              onUpdate={props.onUpdate}
              value={field.value}
            />
          );
        }}
        control={props.control}
        name={props.name}
      />
    </FormControl>
  );
});
