import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Text } from "@/components/ui/text";
import { Grid } from "@/styled-system/jsx";

import { ProgressCircle, ProgressHorizontal } from ".";

const meta = {
  title: "UI/Progress",
  component: ProgressHorizontal,
  args: {
    value: 64,
    size: "md",
    children: "Upload progress",
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<typeof ProgressHorizontal>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Horizontal: Story = {
  render: (args) => (
    <ProgressHorizontal {...args} maxW="md">
      {args.children}
    </ProgressHorizontal>
  ),
};

export const Circle: Story = {
  render: (args) => <ProgressCircle {...args}>{args.children}</ProgressCircle>,
};

export const Sizes: Story = {
  render: () => (
    <Grid
      alignItems="center"
      columnGap="6"
      gridTemplateColumns="3rem minmax(0, 1fr) 4rem"
      maxW="2xl"
      rowGap="5"
      w="full"
    >
      <Text color="text.subtle" fontWeight="semibold" size="xs">
        Size
      </Text>
      <Text color="text.subtle" fontWeight="semibold" size="xs">
        Linear
      </Text>
      <Text color="text.subtle" fontWeight="semibold" size="xs">
        Circular
      </Text>

      <Text fontWeight="semibold" size="sm">
        SM
      </Text>
      <ProgressHorizontal size="sm" value={64} />
      <ProgressCircle size="sm" value={64} />

      <Text fontWeight="semibold" size="sm">
        MD
      </Text>
      <ProgressHorizontal size="md" value={64} />
      <ProgressCircle size="md" value={64} />

      <Text fontWeight="semibold" size="sm">
        LG
      </Text>
      <ProgressHorizontal size="lg" value={64} />
      <ProgressCircle size="lg" value={64} />
    </Grid>
  ),
};
