import {
  type ColorPickerValueChangeDetails,
  parseColor,
} from "@ark-ui/react/color-picker";
import {
  Controller,
  type ControllerProps,
  type FieldPathByValue,
  type FieldValues,
} from "react-hook-form";

import { HStack, Stack } from "@/styled-system/jsx";

import { IconButton } from "../icon-button";
import { ColourPipetteIcon } from "../icons/Colour";
import { Input } from "../input";

import * as ColorPicker from "./ColorPicker.internal";

const presets = [
  "hsl(10, 81%, 59%)",
  "hsl(60, 81%, 59%)",
  "hsl(100, 81%, 59%)",
  "hsl(175, 81%, 59%)",
  "hsl(190, 81%, 59%)",
  "hsl(205, 81%, 59%)",
  "hsl(220, 81%, 59%)",
  "hsl(250, 81%, 59%)",
  "hsl(280, 81%, 59%)",
  "hsl(350, 81%, 59%)",
];

type StringFieldName<TFieldValues extends FieldValues> = FieldPathByValue<
  TFieldValues,
  string | undefined
>;

export type ColorPickerFieldProps<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
> = Omit<ControllerProps<TFieldValues, TName>, "render"> & {
  ariaLabel: string;
};

export function ColorPickerField<
  TFieldValues extends FieldValues,
  TName extends StringFieldName<TFieldValues> = StringFieldName<TFieldValues>,
>({
  ariaLabel,
  ...controllerProps
}: ColorPickerFieldProps<TFieldValues, TName>) {
  return (
    <Controller
      {...controllerProps}
      render={({ field, fieldState }) => (
        <ColorPicker.Root
          name={field.name}
          value={safeParseColour(field.value)}
          onValueChange={({ valueAsString }: ColorPickerValueChangeDetails) =>
            field.onChange(valueAsString)
          }
          onOpenChange={({ open }) => {
            if (!open) field.onBlur();
          }}
          disabled={field.disabled}
          invalid={fieldState.invalid}
        >
          <ColorPicker.Label srOnly>{ariaLabel}</ColorPicker.Label>
          <ColorPicker.Context>
            {(api) => (
              <>
                <ColorPicker.Control>
                  <ColorPicker.ChannelInput channel="hex" asChild>
                    <Input />
                  </ColorPicker.ChannelInput>
                  <ColorPicker.Trigger asChild>
                    <IconButton variant="outline" aria-label="Choose colour">
                      <ColorPicker.Swatch value={api.value} />
                    </IconButton>
                  </ColorPicker.Trigger>
                </ColorPicker.Control>
                <ColorPicker.Positioner>
                  <ColorPicker.Content>
                    <Stack gap="3">
                      <ColorPicker.Area>
                        <ColorPicker.AreaBackground />
                        <ColorPicker.AreaThumb />
                      </ColorPicker.Area>
                      <HStack gap="3">
                        <ColorPicker.EyeDropperTrigger asChild>
                          <IconButton
                            size="sm"
                            variant="outline"
                            aria-label="Pick a colour"
                          >
                            <ColourPipetteIcon />
                          </IconButton>
                        </ColorPicker.EyeDropperTrigger>
                        <Stack gap="2" flex="1">
                          <ColorPicker.ChannelSlider channel="hue">
                            <ColorPicker.ChannelSliderTrack />
                            <ColorPicker.ChannelSliderThumb />
                          </ColorPicker.ChannelSlider>
                        </Stack>
                      </HStack>
                      <HStack>
                        <ColorPicker.ChannelInput channel="hex" asChild>
                          <Input h="6" minW="7" px="1.5" fontSize="xs" />
                        </ColorPicker.ChannelInput>
                      </HStack>
                      <Stack gap="1.5">
                        <ColorPicker.SwatchGroup>
                          {presets.map((color) => (
                            <ColorPicker.SwatchTrigger
                              key={color}
                              value={color}
                            >
                              <ColorPicker.Swatch value={color} />
                            </ColorPicker.SwatchTrigger>
                          ))}
                        </ColorPicker.SwatchGroup>
                      </Stack>
                    </Stack>
                  </ColorPicker.Content>
                </ColorPicker.Positioner>
              </>
            )}
          </ColorPicker.Context>
          <ColorPicker.HiddenInput ref={field.ref} />
        </ColorPicker.Root>
      )}
    />
  );
}

function safeParseColour(value: unknown) {
  if (typeof value === "string") {
    try {
      return parseColor(value);
    } catch {
      // Fall through to the neutral valid value below.
    }
  }

  return parseColor("green");
}
