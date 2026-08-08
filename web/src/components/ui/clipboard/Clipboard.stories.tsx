import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "../button";
import { Input } from "../input";

import * as Clipboard from ".";

const meta = {
  title: "Components/Actions/Clipboard",
  component: Clipboard.Root,
} satisfies Meta<typeof Clipboard.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Clipboard.Root value="https://storyden.org" maxW="md">
      <Clipboard.Label>Invite link</Clipboard.Label>
      <Clipboard.Control>
        <Clipboard.Input asChild>
          <Input readOnly />
        </Clipboard.Input>
        <Clipboard.Trigger asChild>
          <Button variant="outline">Copy</Button>
        </Clipboard.Trigger>
      </Clipboard.Control>
    </Clipboard.Root>
  ),
};
