import { CSSProperties, Ref, forwardRef } from "react";
import { Control, Controller } from "react-hook-form";

import { FormControl } from "@/components/ui/form/FormControl";
import { DragHandleIcon } from "@/components/ui/icons/DragHandle";
import { styled } from "@/styled-system/jsx";

import { C, L, Props, useColourInput } from "./useColourInput";

export function ColourInput(props: Props) {
  const { onPointerDown, onPointerUp, hue, ref, angle, value, grabbing } =
    useColourInput(props);

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
      borderColor="border.muted"
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
        boxShadow="lg"
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
            <DragHandleIcon color="fg.muted" />
          </styled.div>
        </styled.output>
      </styled.div>
    </styled.div>
  );
}

type FieldProps = {
  control: Control<any>;
  name: string;
  defaultValue: string;
  onUpdate: (v: string) => void;
};

const _ColourField = (props: FieldProps, ref: Ref<HTMLDivElement>) => {
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
};

export const ColourField = forwardRef(_ColourField);
