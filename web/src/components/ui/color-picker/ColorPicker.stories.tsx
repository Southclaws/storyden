import { parseColor } from "@ark-ui/react/color-picker";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { HStack, Stack } from "@/styled-system/jsx";

import { IconButton } from "../icon-button";
import { ColourPipetteIcon } from "../icons/Colour";
import { Input } from "../input";

import * as ColorPicker from ".";

const presets = [
  "hsl(10, 81%, 59%)",
  "hsl(60, 81%, 59%)",
  "hsl(100, 81%, 59%)",
  "hsl(175, 81%, 59%)",
  "hsl(220, 81%, 59%)",
  "hsl(280, 81%, 59%)",
  "hsl(350, 81%, 59%)",
];

const meta = {
  title: "UI/Color Picker",
  component: ColorPicker.Root,
} satisfies Meta<typeof ColorPicker.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <ColorPicker.Root defaultValue={parseColor("hsl(175, 81%, 39%)")}>
      <ColorPicker.Context>
        {(api) => (
          <>
            <ColorPicker.Control>
              <ColorPicker.ChannelInput channel="hex" asChild>
                <Input maxW="44" />
              </ColorPicker.ChannelInput>
              <ColorPicker.Trigger asChild>
                <IconButton variant="outline" aria-label="Open color picker">
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
                  <ColorPicker.SwatchGroup>
                    {presets.map((color) => (
                      <ColorPicker.SwatchTrigger key={color} value={color}>
                        <ColorPicker.Swatch value={color} />
                      </ColorPicker.SwatchTrigger>
                    ))}
                  </ColorPicker.SwatchGroup>
                </Stack>
              </ColorPicker.Content>
            </ColorPicker.Positioner>
          </>
        )}
      </ColorPicker.Context>
      <ColorPicker.HiddenInput />
    </ColorPicker.Root>
  ),
};
