import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import type { ThemeAsset } from "@/api/openapi-schema";

import { ThemeSettingsView } from "./ThemeSettingsView";

const stylesheet = asset("theme-css", "text/css");
const script = asset("theme-js", "application/javascript");

const meta = {
  title: "Compositions/Admin/Custom Theme",
  component: ThemeSettingsView,
  parameters: { layout: "fullscreen" },
  decorators: [
    (Story) => (
      <div style={{ padding: "2rem", maxWidth: "72rem" }}>
        <Story />
      </div>
    ),
  ],
  args: {
    assets: [],
    active: false,
    runtimeDisabled: false,
    editingEnabled: false,
    onEnableEditing: () => undefined,
    onExitEditing: () => undefined,
    onDisable: () => undefined,
  },
} satisfies Meta<typeof ThemeSettingsView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Empty: Story = {};
export const Active: Story = {
  args: { assets: [stylesheet], active: true },
};
export const EditingEnabled: Story = {
  args: { assets: [stylesheet], active: true, editingEnabled: true },
};
export const RuntimeDisabled: Story = {
  args: { assets: [stylesheet], active: true, runtimeDisabled: true },
};
export const ValidationError: Story = {
  args: { error: "The active theme exceeds the 5 MiB total limit." },
};
export const DangerousScriptWarning: Story = {
  args: { assets: [stylesheet, script], active: true },
};

function asset(id: string, mimeType: ThemeAsset["mime_type"]): ThemeAsset {
  return {
    id: id.padEnd(20, "0"),
    filename: id,
    path: `/api/info/theme/assets/${id}`,
    mime_type: mimeType,
    size: 768,
    integrity: "sha256-aW5kZXBlbmRlbnQtZXhwZWN0ZWQtZGlnZXN0",
  };
}
