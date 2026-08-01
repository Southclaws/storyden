import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { styled } from "@/styled-system/jsx";

import { ScrollToTop } from ".";

const meta = {
  title: "UI/Scroll To Top",
  component: ScrollToTop,
} satisfies Meta<typeof ScrollToTop>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <styled.div style={{ minHeight: "120vh" }} display="flex" alignItems="end">
      <ScrollToTop>Back to top</ScrollToTop>
    </styled.div>
  ),
};
