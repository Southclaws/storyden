import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { HStack, LStack, styled } from "@/styled-system/jsx";

import { LikeButton } from ".";

const meta = {
  title: "Components/Actions/Like Button",
  component: LikeButton,
  parameters: {
    docs: {
      description: {
        component:
          "Use LikeButton to toggle a like and display its count. Both sizes keep the count inside the button so it shares the control background over images. Small is for repeated thread cards: a heart and count. Large is for the thread byline and adds the Like label. Use IconButton for unrelated icon actions and ReactList for emoji reactions. Guests retain a disabled control and count; pending controls preserve focus and reject repeat clicks. The component is controlled and does not fetch or mutate data.",
      },
    },
  },
  args: {
    size: "small",
    count: 2,
    liked: false,
    enabled: true,
    pending: false,
  },
  argTypes: {
    size: { control: "select", options: ["small", "large"] },
  },
} satisfies Meta<typeof LikeButton>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Small: Story = {};

export const Large: Story = {
  args: { size: "large" },
};

export const OverImage: Story = {
  render: (args) => (
    <styled.div
      backgroundImage={`url(${new URL("../../../stories/assets/mountains.png", import.meta.url).href})`}
      backgroundSize="cover"
      backgroundPosition="center"
      borderRadius="panel"
      padding="6"
    >
      <HStack gap="4" flexWrap="wrap">
        <LikeButton {...args} liked />
        <LikeButton {...args} size="large" liked />
      </HStack>
    </styled.div>
  ),
};

export const States: Story = {
  render: (args) => (
    <LStack gap="4" alignItems="start">
      {(["small", "large"] as const).map((size) => (
        <HStack key={size} gap="4" flexWrap="wrap">
          <LikeButton {...args} size={size} />
          <LikeButton {...args} size={size} liked />
          <LikeButton {...args} size={size} liked pending />
          <LikeButton {...args} size={size} enabled={false} />
        </HStack>
      ))}
    </LStack>
  ),
};

export const Counts: Story = {
  render: (args) => (
    <LStack gap="4" alignItems="start">
      {(["small", "large"] as const).map((size) => (
        <HStack key={size} gap="4" flexWrap="wrap">
          <LikeButton {...args} size={size} count={0} />
          <LikeButton {...args} size={size} count={1} />
          <LikeButton {...args} size={size} count={12345} />
        </HStack>
      ))}
    </LStack>
  ),
};

export const Interactive: Story = {
  args: { size: "large" },
  render: function Interactive(args) {
    const [liked, setLiked] = useState(args.liked);
    const [count, setCount] = useState(args.count);

    function toggle() {
      if (liked) {
        setCount(count - 1);
      } else {
        setCount(count + 1);
      }

      setLiked(!liked);
    }

    return (
      <LikeButton {...args} liked={liked} count={count} onClick={toggle} />
    );
  },
};
