import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { BackAction } from "@/components/site/Action/Back";
import { Button } from "@/components/ui/button";
import { LinkButton } from "@/components/ui/link-button";
import { LStack } from "@/styled-system/jsx";

import { PageHeader } from ".";

const meta = {
  title: "Compositions/Page headers",
  component: PageHeader,
  parameters: {
    layout: "padded",
  },
} satisfies Meta<typeof PageHeader>;

export default meta;

type Story = StoryObj<typeof meta>;

export const HeadingOnly: Story = {
  args: {
    title: "Members",
  },
};

export const WithDescriptionAndAction: Story = {
  args: {
    title: "Discussion categories",
    description: "There are 4 categories available to start discussions.",
    actions: <Button>Create</Button>,
  },
};

export const BackNavigation: Story = {
  args: {
    title: "Library",
    back: <BackAction fallbackHref="/" />,
    actions: <Button>New page</Button>,
  },
};

export const BreadcrumbNavigation: Story = {
  args: {
    title: "",
  },
  render: () => (
    <LStack gap="2">
      <PageHeader
        title="API reference"
        navigation={
          <>
            <LinkButton href="/l" variant="ghost">
              Library
            </LinkButton>
            <LinkButton href="/l/documentation" variant="ghost">
              Documentation
            </LinkButton>
          </>
        }
        actions={<Button>Edit</Button>}
      />
    </LStack>
  ),
};
