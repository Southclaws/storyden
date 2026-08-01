import { createListCollection } from "@ark-ui/react";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "@/components/ui/button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { SelectIcon } from "@/components/ui/icons/Select";
import { Input } from "@/components/ui/input";
import * as Select from "@/components/ui/select";
import { styled } from "@/styled-system/jsx";

const contentTypes = [
  { label: "Discussion", value: "discussion" },
  { label: "Library page", value: "library" },
  { label: "Collection", value: "collection" },
];

const collection = createListCollection({ items: contentTypes });

const sizes = ["sm", "md", "lg"] as const;

const meta = {
  title: "Compositions/Forms/Control Sizing",
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;
type ControlSize = (typeof sizes)[number];

function ExampleSelect({ size }: { size?: ControlSize }) {
  return (
    <Select.Root collection={collection} size={size}>
      <Select.Control>
        <Select.Trigger>
          <Select.ValueText placeholder="Choose type" />
          <SelectIcon />
        </Select.Trigger>
      </Select.Control>
      <Select.Positioner>
        <Select.Content>
          {contentTypes.map((item) => (
            <Select.Item key={item.value} item={item}>
              <Select.ItemText>{item.label}</Select.ItemText>
              <Select.ItemIndicator>
                <CheckIcon />
              </Select.ItemIndicator>
            </Select.Item>
          ))}
        </Select.Content>
      </Select.Positioner>
      <Select.HiddenSelect />
    </Select.Root>
  );
}

function ControlRow({ label, size }: { label: string; size: ControlSize }) {
  return (
    <styled.div
      alignItems="start"
      borderBottomColor="border.subtle"
      borderBottomWidth="thin"
      display="grid"
      gap="3"
      gridTemplateColumns={{
        base: "1fr",
        md: "8rem 10rem minmax(12rem, 1fr) minmax(12rem, 1fr)",
      }}
      paddingY="4"
    >
      <styled.div display="flex" flexDirection="column" gap="1">
        <styled.strong
          color="content.default"
          fontSize="sm"
          fontWeight="medium"
        >
          {label}
        </styled.strong>
        <styled.code color="content.muted" fontFamily="mono" fontSize="xs">
          {`size="${size}"`}
        </styled.code>
      </styled.div>
      <Button size={size} variant="subtle">
        Save
      </Button>
      <Input placeholder="storyden" size={size} />
      <ExampleSelect size={size} />
    </styled.div>
  );
}

export const ControlSizing: Story = {
  render: () => (
    <styled.main
      boxSizing="border-box"
      color="content.default"
      display="flex"
      flexDirection="column"
      gap="8"
      marginX="auto"
      maxW="6xl"
      padding={{ base: "4", md: "8" }}
      width="full"
    >
      <styled.header display="flex" flexDirection="column" gap="2" maxW="3xl">
        <styled.h1 fontSize="2xl" fontWeight="semibold" lineHeight="tight">
          Form control sizing
        </styled.h1>
        <styled.p color="content.subtle" fontSize="md" lineHeight="relaxed">
          Compare Button, Input, and Select sizing in the same form row while
          tuning dense product defaults.
        </styled.p>
      </styled.header>

      <styled.section display="flex" flexDirection="column">
        <styled.div
          alignItems="end"
          borderBottomColor="border.subtle"
          borderBottomWidth="thin"
          color="content.subtle"
          display="grid"
          fontSize="sm"
          fontWeight="medium"
          gap="3"
          gridTemplateColumns={{
            base: "1fr",
            md: "8rem 10rem minmax(12rem, 1fr) minmax(12rem, 1fr)",
          }}
          paddingBottom="2"
        >
          <styled.div />
          <styled.div>Button</styled.div>
          <styled.div>Input</styled.div>
          <styled.div>Select</styled.div>
        </styled.div>
        {sizes.map((size) => (
          <ControlRow key={size} label={size.toUpperCase()} size={size} />
        ))}
      </styled.section>
    </styled.main>
  ),
};
