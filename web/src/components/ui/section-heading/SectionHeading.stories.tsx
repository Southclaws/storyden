import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { SectionHeading } from ".";

const meta = {
  title: "Components/Typography/Section Heading",
  component: SectionHeading,
  parameters: {
    docs: {
      description: {
        component:
          "Introduces a meaningful subdivision within a product screen. SectionHeading always renders `h2` and owns the compact product hierarchy, so do not recreate it with `styled.h2` and local typography props. Do not use it merely to label a Card; Surface titles belong to the Surface recipe.",
      },
    },
  },
  args: {
    children: "Client IP strategy",
  },
} satisfies Meta<typeof SectionHeading>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Strong: Story = {
  args: {
    children: "Robot actions",
    emphasis: "strong",
  },
};
