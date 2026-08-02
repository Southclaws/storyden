export type TokenExample = {
  name: string;
  token: string;
  usage: string;
};

export type StatusTokenExample = {
  name: string;
  surface: string;
  content: string;
  border: string;
};

export type LayerRoleTokenExample = {
  name: string;
  usage: string;
  background: string;
  elevated?: boolean;
};

export const layerRoleTokens: LayerRoleTokenExample[] = [
  {
    name: "Canvas",
    usage: "Root application background and quiet page-level content.",
    background: "colors.background.canvas",
  },
  {
    name: "Surface",
    usage: "Independent objects such as posts, cards, panels, and reports.",
    background: "colors.background.surface",
  },
  {
    name: "Inset",
    usage: "Embedded or supporting content within the current page context.",
    background: "colors.background.inset",
  },
  {
    name: "Overlay",
    usage:
      "Temporary content above document flow, including menus and dialogs.",
    background: "colors.background.overlay",
    elevated: true,
  },
];

export const textTokens: TokenExample[] = [
  {
    name: "Default text",
    token: "colors.text.default",
    usage: "Primary product copy and labels.",
  },
  {
    name: "Muted text",
    token: "colors.text.muted",
    usage: "Secondary descriptions and metadata.",
  },
  {
    name: "Subtle text",
    token: "colors.text.subtle",
    usage: "Low-priority supporting information.",
  },
  {
    name: "Disabled text",
    token: "colors.text.disabled",
    usage: "Unavailable actions and values.",
  },
  {
    name: "Inverse text",
    token: "colors.text.inverse",
    usage: "Text placed on a strong filled control.",
  },
];

export const borderTokens: TokenExample[] = [
  {
    name: "Default border",
    token: "colors.border.default",
    usage: "Ordinary separation between adjacent regions.",
  },
  {
    name: "Muted border",
    token: "colors.border.muted",
    usage: "Lower-emphasis control and supporting-region boundaries.",
  },
  {
    name: "Strong border",
    token: "colors.border.strong",
    usage: "Overlays and boundaries that need clearer separation.",
  },
  {
    name: "Disabled border",
    token: "colors.border.disabled",
    usage: "Unavailable controls.",
  },
];

export const stateBackgroundTokens: TokenExample[] = [
  {
    name: "Control",
    token: "colors.background.control",
    usage: "Persistent fill for inputs and other controls.",
  },
  {
    name: "Control hover",
    token: "colors.background.controlHover",
    usage: "Neutral hover feedback for interactive controls and rows.",
  },
  {
    name: "Control subtle",
    token: "colors.background.controlSubtle",
    usage: "Low-emphasis filled controls such as subtle buttons.",
  },
  {
    name: "Control track",
    token: "colors.background.controlTrack",
    usage: "Tracks for switches, sliders, and progress indicators.",
  },
  {
    name: "Selection",
    token: "colors.selection.background",
    usage: "Neutral selected state without semantic accent meaning.",
  },
];

export const statusTokens: StatusTokenExample[] = [
  {
    name: "danger",
    surface: "colors.status.danger.surface",
    content: "colors.status.danger.content",
    border: "colors.status.danger.border",
  },
  {
    name: "success",
    surface: "colors.status.success.surface",
    content: "colors.status.success.content",
    border: "colors.status.success.border",
  },
  {
    name: "warning",
    surface: "colors.status.warning.surface",
    content: "colors.status.warning.content",
    border: "colors.status.warning.border",
  },
  {
    name: "info",
    surface: "colors.status.info.surface",
    content: "colors.status.info.content",
    border: "colors.status.info.border",
  },
];

export const visibilityTokens: StatusTokenExample[] = [
  {
    name: "published",
    surface: "colors.visibility.published.bg",
    content: "colors.visibility.published.fg",
    border: "colors.visibility.published.border",
  },
  {
    name: "draft",
    surface: "colors.visibility.draft.bg",
    content: "colors.visibility.draft.fg",
    border: "colors.visibility.draft.border",
  },
  {
    name: "review",
    surface: "colors.visibility.review.bg",
    content: "colors.visibility.review.fg",
    border: "colors.visibility.review.border",
  },
  {
    name: "unlisted",
    surface: "colors.visibility.unlisted.bg",
    content: "colors.visibility.unlisted.fg",
    border: "colors.visibility.unlisted.border",
  },
];

export const accentScaleTokens: TokenExample[] = [
  "1",
  "2",
  "3",
  "4",
  "5",
  "6",
  "7",
  "8",
  "9",
  "10",
].map((step) => ({
  name: `Accent ${step}`,
  token: `colors.accent.${step}`,
  usage: "Community-configurable accent ramp.",
}));

export const primitiveColorRamps = [
  "neutral",
  "slate",
  "blue",
  "green",
  "amber",
  "red",
] as const;

export const spacingTokens = [
  "0",
  "0.5",
  "1",
  "1.5",
  "2",
  "2.5",
  "3",
  "3.5",
  "4",
  "4.5",
  "5",
  "5.5",
  "6",
  "7",
  "7.5",
  "8",
  "9",
  "9.5",
  "10",
  "10.5",
  "11",
  "12",
  "14",
  "16",
  "20",
  "24",
  "28",
  "32",
  "36",
  "40",
  "44",
  "48",
  "52",
  "56",
  "60",
  "64",
  "72",
  "80",
  "96",
] as const;

export const layoutSizeTokens: TokenExample[] = [
  {
    name: "Readable",
    token: "sizes.layout.readable",
    usage: "Long-form text measure.",
  },
  {
    name: "Content",
    token: "sizes.layout.content",
    usage: "Default main content width.",
  },
  {
    name: "Wide content",
    token: "sizes.layout.contentWide",
    usage: "Wide composed screens.",
  },
  {
    name: "Form",
    token: "sizes.layout.form",
    usage: "Single-column form width.",
  },
  {
    name: "Sidebar",
    token: "sizes.layout.sidebar",
    usage: "Desktop navigation/sidebar width.",
  },
  {
    name: "Drawer",
    token: "sizes.layout.drawer",
    usage: "Tablet drawer width.",
  },
  {
    name: "Command bar",
    token: "sizes.layout.commandBar",
    usage: "Mobile command bar width.",
  },
];

export const radiusTokens = [
  "radii.none",
  "radii.2xs",
  "radii.xs",
  "radii.sm",
  "radii.md",
  "radii.lg",
  "radii.xl",
  "radii.2xl",
  "radii.3xl",
  "radii.full",
  "radii.control",
  "radii.panel",
  "radii.overlay",
  "radii.pill",
] as const;

export const shadowTokens = [
  "shadows.surface",
  "shadows.floating",
  "shadows.overlay",
] as const;

export const blurTokens = [
  "blurs.sm",
  "blurs.base",
  "blurs.md",
  "blurs.lg",
  "blurs.xl",
  "blurs.subtle",
  "blurs.frosted",
] as const;

export const fontFamilyTokens: TokenExample[] = [
  {
    name: "Body",
    token: "fonts.body",
    usage: "Product UI text.",
  },
  {
    name: "Heading",
    token: "fonts.heading",
    usage: "Headings and display text.",
  },
  {
    name: "Mono",
    token: "fonts.mono",
    usage: "Code, IDs, and token labels.",
  },
];

export const fontSizeTokens = [
  "2xs",
  "xs",
  "sm",
  "md",
  "lg",
  "xl",
  "2xl",
] as const;

export const breakpointTokens = [
  {
    name: "sm",
    value: 640,
    usage: "Small responsive adjustments.",
  },
  {
    name: "md",
    value: 768,
    usage: "Tablet and mobile navigation boundary.",
  },
  {
    name: "lg",
    value: 1024,
    usage: "Desktop shell and content layout.",
  },
  {
    name: "xl",
    value: 1280,
    usage: "Wide desktop gutters.",
  },
  {
    name: "2xl",
    value: 1536,
    usage: "Very wide viewport refinement.",
  },
] as const;
