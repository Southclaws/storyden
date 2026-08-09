import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { BackAction } from "@/components/site/Action/Back";
import { Breadcrumbs } from "@/components/ui/breadcrumbs";
import { Button } from "@/components/ui/button";

import { PageHeader } from ".";

const meta = {
  title: "Compositions/Page Layout/Page Header",
  component: PageHeader,
  parameters: {
    layout: "padded",
    docs: {
      description: {
        component:
          "The canonical top-of-page composition for a page heading, supporting text, return/breadcrumb navigation, and a small action group. Omit it only when the content makes the page identity genuinely obvious. Breadcrumbs describe hierarchy; BackAction describes a return path.",
      },
    },
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
    <PageHeader
      title="Web Development"
      navigation={
        <Breadcrumbs
          index={{ label: "Discussion", href: "/d" }}
          crumbs={[
            {
              label: "Web Development",
              href: "/d/web-development",
            },
          ]}
        />
      }
    />
  ),
};
