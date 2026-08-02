import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Text } from "@/components/ui/text";
import { Grid } from "@/styled-system/jsx";

import { NumberInput } from ".";

const meta = {
  title: "UI/Number Input",
  component: NumberInput,
  args: {
    defaultValue: "12",
    min: 0,
    max: 100,
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof NumberInput>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Playground: Story = {
  render: (args) => <NumberInput {...args} maxW="sm" />,
};

export const Sizes: Story = {
  render: () => (
    <Grid
      alignItems="center"
      columnGap="4"
      gridTemplateColumns="3rem minmax(0, 1fr)"
      maxW="sm"
      rowGap="3"
      w="full"
    >
      <Text fontWeight="semibold" variant="supporting">
        SM
      </Text>
      <NumberInput
        size="sm"
        defaultValue="12"
        inputProps={{ "aria-label": "Number input sm" }}
      />

      <Text fontWeight="semibold" variant="supporting">
        MD
      </Text>
      <NumberInput
        size="md"
        defaultValue="12"
        inputProps={{ "aria-label": "Number input md" }}
      />

      <Text fontWeight="semibold" variant="supporting">
        LG
      </Text>
      <NumberInput
        size="lg"
        defaultValue="12"
        inputProps={{ "aria-label": "Number input lg" }}
      />
    </Grid>
  ),
};

export const Scrubber: Story = {
  render: () => (
    <NumberInput
      scrubber
      defaultValue="36"
      maxW="sm"
      inputProps={{ "aria-label": "Scrubbable number input" }}
    />
  ),
};
