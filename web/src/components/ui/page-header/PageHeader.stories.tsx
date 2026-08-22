import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { BackAction } from "@/components/site/Action/Back";
import { Badge } from "@/components/ui/badge";
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
          "The canonical top-of-page composition for a page heading, optional supporting text, return or breadcrumb navigation, and a small action group. Use it instead of building a feature-local header row from Flex, Grid, or Stack. Actions remain aligned with the first title line while long identity wraps into the available space. Supporting text must add decision-relevant context; omit it when the title and page content already explain the task. Breadcrumbs describe hierarchy; BackAction describes a return path.",
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

export const FullPageForm: Story = {
  args: {
    title: "New Trail",
    back: <BackAction fallbackHref="/robots/trails" />,
  },
  parameters: {
    docs: {
      description: {
        story:
          "Full-page forms use the back action as their exit path. Omit a description when the fields and section headings provide the necessary guidance.",
      },
    },
  },
};

export const LongIdentity: Story = {
  args: {
    title: "W".repeat(120),
    description: "W".repeat(240),
    badge: (
      <Badge colorPalette="green" size="sm">
        Active
      </Badge>
    ),
    back: <BackAction fallbackHref="/" />,
    actions: <Button variant="solid">Primary action</Button>,
  },
  parameters: {
    docs: {
      description: {
        story:
          "Long unbroken titles and descriptions wrap inside the available heading space. The return action and primary action stay aligned with the first title line at every width.",
      },
    },
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
