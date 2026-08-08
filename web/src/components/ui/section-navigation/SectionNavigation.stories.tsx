import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Box } from "@/styled-system/jsx";

import { SectionNavigation } from ".";

const meta = {
  title: "Components/Navigation/Section Navigation",
  component: SectionNavigation,
  parameters: {
    layout: "padded",
  },
  decorators: [
    (Story) => (
      <Box width="full" maxWidth="[24rem]">
        <Story />
      </Box>
    ),
  ],
} satisfies Meta<typeof SectionNavigation>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Settings: Story = {
  args: {
    activeHref: "/settings/authentication",
    defaultOpen: true,
    label: "Settings sections",
    groups: [
      {
        items: [
          { href: "/settings/interface", label: "Interface" },
          { href: "/settings/authentication", label: "Authentication" },
          { href: "/settings/email", label: "Email" },
          { href: "/settings/access-keys", label: "Access keys" },
        ],
      },
    ],
  },
};

export const GroupedAdmin: Story = {
  args: {
    activeHref: "/admin/authentication",
    defaultOpen: true,
    label: "Admin sections",
    groups: [
      {
        label: "Appearance",
        items: [
          { href: "/admin/brand", label: "Brand" },
          { href: "/admin/interface", label: "Interface" },
        ],
      },
      {
        label: "Security",
        items: [
          { href: "/admin/authentication", label: "Authentication" },
          { href: "/admin/access-keys", label: "Access keys" },
          { href: "/admin/oauth", label: "OAuth" },
        ],
      },
      {
        label: "Operations",
        items: [
          { href: "/admin/system", label: "System" },
          { href: "/admin/audit-log", label: "Audit log" },
          { href: "/admin/email", label: "Email" },
        ],
      },
    ],
  },
};
