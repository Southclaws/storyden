import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Box } from "@/styled-system/jsx";

import { DenDoctor } from "./DenDoctor";

const meta = {
  title: "Compositions/Admin/Den Doctor",
  component: DenDoctor,
  decorators: [
    (Story) => (
      <Box p="4" maxW="3xl">
        <Story />
      </Box>
    ),
  ],
  parameters: {
    docs: {
      description: {
        component:
          "Local deployment-address guidance for a fresh installation and System settings. Reuses Alert for persistent problems and native details for operator configuration. Does not probe servers or diagnose DNS, TLS certificates, or CORS responses.",
      },
    },
  },
  args: {
    apiAddress: "https://api.example.com",
    webAddress: "https://example.com",
    browserOrigin: "https://example.com",
  },
} satisfies Meta<typeof DenDoctor>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Healthy: Story = { args: { showHealthy: true } };
export const FreshHealthy: Story = {};
export const MixedSchemes: Story = {
  args: { apiAddress: "http://api.example.com" },
};
export const SiblingHosts: Story = {
  args: {
    webAddress: "https://www.example.com",
    browserOrigin: "https://www.example.com",
  },
};
export const LocalhostAPI: Story = {
  args: { apiAddress: "http://localhost:8000" },
};
export const DifferentBrowserAddress: Story = {
  args: { browserOrigin: "https://www.example.com" },
};
export const WebsiteSubpath: Story = {
  args: { webAddress: "https://example.com/community" },
};
export const APIProxyPrefix: Story = {
  args: { apiAddress: "https://api.example.com/backend" },
};
