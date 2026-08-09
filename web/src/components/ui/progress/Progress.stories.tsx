import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Text } from "@/components/ui/text";
import { Grid } from "@/styled-system/jsx";

import { ProgressCircle, ProgressHorizontal } from ".";

const meta = {
  title: "Components/Feedback/Progress",
  component: ProgressHorizontal,
  parameters: {
    docs: {
      description: {
        component:
          "Shows measurable completion for an operation with known progress. Use the horizontal form when comparison or a label matters and the circular form in compact chrome. Use Spinner for indeterminate waiting and a skeleton when preserving the shape of loading content.",
      },
    },
  },
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
      <Text fontWeight="semibold" variant="metadata">
        Size
      </Text>
      <Text fontWeight="semibold" variant="metadata">
        Linear
      </Text>
      <Text fontWeight="semibold" variant="metadata">
        Circular
      </Text>

      <Text fontWeight="semibold" variant="supporting">
        SM
      </Text>
      <ProgressHorizontal size="sm" value={64} />
      <ProgressCircle size="sm" value={64} />

      <Text fontWeight="semibold" variant="supporting">
        MD
      </Text>
      <ProgressHorizontal size="md" value={64} />
      <ProgressCircle size="md" value={64} />

      <Text fontWeight="semibold" variant="supporting">
        LG
      </Text>
      <ProgressHorizontal size="lg" value={64} />
      <ProgressCircle size="lg" value={64} />
    </Grid>
  ),
};
