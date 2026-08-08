import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { css } from "@/styled-system/css";
import { HStack, styled } from "@/styled-system/jsx";

import { RobotActivityIcon } from "./RobotActivity";
import { RunningAnimatedIcon } from "./RunningAnimatedIcon";
import { SiteIcon } from "./Site";

const meta = {
  title: "Foundations/Icons/Special",
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

const iconPreviewClassName = css({
  "& svg, & img": {
    width: "[100%]",
    height: "[100%]",
    maxWidth: "[100%]",
    maxHeight: "[100%]",
  },
});

export const Special: Story = {
  render: () => (
    <HStack gap="4" alignItems="center">
      <IconTile name="SiteIcon">
        <SiteIcon width="5" height="5" />
      </IconTile>
      <IconTile name="RobotActivityIcon">
        <RobotActivityIcon size={20} />
      </IconTile>
      <IconTile name="RunningAnimatedIcon">
        <RunningAnimatedIcon size={20} />
      </IconTile>
    </HStack>
  ),
};

function IconTile({
  name,
  children,
}: {
  name: string;
  children: React.ReactNode;
}) {
  return (
    <styled.div
      display="flex"
      alignItems="center"
      gap="2"
      borderWidth="thin"
      borderColor="border.muted"
      borderRadius="md"
      p="2"
    >
      <styled.span
        color="text.default"
        flexShrink="0"
        display="inline-flex"
        alignItems="center"
        justifyContent="center"
        width="5"
        height="5"
        overflow="hidden"
        className={iconPreviewClassName}
      >
        {children}
      </styled.span>
      <styled.span color="text.subtle" fontSize="xs">
        {name}
      </styled.span>
    </styled.div>
  );
}
