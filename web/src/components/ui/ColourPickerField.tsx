import {
  ColorPickerValueChangeDetails,
  parseColor,
} from "@ark-ui/react/color-picker";
import { Control, Controller, FieldValues, Path } from "react-hook-form";

import * as ColorPicker from "@/components/ui/color-picker";
import { HStack, Stack } from "@/styled-system/jsx";

import { IconButton } from "./icon-button";
import { ColourPipetteIcon } from "./icons/Colour";
import { Input } from "./input";

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

export type Props<T extends FieldValues> = {
  control: Control<T>;
  name: Path<T>;
};

export function ColourPickerField<T extends FieldValues>(props: Props<T>) {
  return (
    <Controller<T>
      control={props.control}
      name={props.name}
      render={({ field }) => {
        function handleChange(d: ColorPickerValueChangeDetails) {
          field.onChange(d.valueAsString);
        }

        const value = safeParseColour(field.value);

        return (
          <ColorPicker.Root value={value} onValueChange={handleChange}>
            <ColorPicker.Context>
              {(api) => (
                <>
                  <ColorPicker.Control>
                    <ColorPicker.ChannelInput channel="hex" asChild>
                      <Input />
                    </ColorPicker.ChannelInput>
                    <ColorPicker.Trigger asChild>
                      <IconButton variant="outline">
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
                              size="xs"
                              variant="outline"
                              aria-label="Pick a color"
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
                            <Input size="2xs" />
                          </ColorPicker.ChannelInput>
                        </HStack>
                        <Stack gap="1.5">
                          <ColorPicker.SwatchGroup>
                            {presets.map((color, id) => (
                              <ColorPicker.SwatchTrigger key={id} value={color}>
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
            <ColorPicker.HiddenInput />
          </ColorPicker.Root>
        );
      }}
    />
  );
}

function safeParseColour(c: string) {
  try {
    return parseColor(c);
  } catch (_) {
    return parseColor("green");
  }
}
