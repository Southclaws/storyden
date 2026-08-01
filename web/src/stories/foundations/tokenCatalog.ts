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

export const surfaceTokens: TokenExample[] = [
  {
    name: "Canvas",
    token: "colors.canvas",
    usage: "Viewport and app background.",
  },
  {
    name: "Surface",
    token: "colors.surface.default",
    usage: "Default panels, cards, and controls.",
  },
  {
    name: "Subtle surface",
    token: "colors.surface.subtle",
    usage: "Secondary surface contrast.",
  },
  {
    name: "Muted surface",
    token: "colors.surface.muted",
    usage: "Hover, quiet selected areas, and nested surfaces.",
  },
  {
    name: "Elevated surface",
    token: "colors.surface.elevated",
    usage: "Raised panels and overlays.",
  },
  {
    name: "Frosted surface",
    token: "colors.surface.frosted",
    usage: "Blurred navigation and translucent panels.",
  },
  {
    name: "Selected surface",
    token: "colors.surface.selected",
    usage: "Selected rows, active nav, and checked cards.",
  },
  {
    name: "Disabled surface",
    token: "colors.surface.disabled",
    usage: "Unavailable component backgrounds.",
  },
];

export const contentTokens: TokenExample[] = [
  {
    name: "Default content",
    token: "colors.content.default",
    usage: "Primary text and icons.",
  },
  {
    name: "Subtle content",
    token: "colors.content.subtle",
    usage: "Secondary labels and metadata.",
  },
  {
    name: "Muted content",
    token: "colors.content.muted",
    usage: "Low-emphasis hints.",
  },
  {
    name: "Disabled content",
    token: "colors.content.disabled",
    usage: "Unavailable text and icons.",
  },
  {
    name: "Accent content",
    token: "colors.content.accent",
    usage: "Community-owned accent foreground.",
  },
];

export const borderTokens: TokenExample[] = [
  {
    name: "Default border",
    token: "colors.border.default",
    usage: "Standard component boundaries.",
  },
  {
    name: "Subtle border",
    token: "colors.border.subtle",
    usage: "Low contrast dividers.",
  },
  {
    name: "Muted border",
    token: "colors.border.muted",
    usage: "Secondary frames and separators.",
  },
  {
    name: "Accent border",
    token: "colors.border.accent",
    usage: "Accent-backed boundaries.",
  },
  {
    name: "Disabled border",
    token: "colors.border.disabled",
    usage: "Inactive component frames.",
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
