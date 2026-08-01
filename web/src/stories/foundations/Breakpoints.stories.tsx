import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { styled } from "@/styled-system/jsx";

import {
  FoundationPage,
  FoundationSection,
  TokenCode,
} from "./FoundationLayout";
import { breakpointTokens } from "./tokenCatalog";

const meta = {
  title: "Foundations/Tokens/Breakpoints",
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

const maxBreakpoint = breakpointTokens[breakpointTokens.length - 1].value;

const containerConditions = [
  {
    name: "containerSmall",
    value: "@container (max-width: 560px)",
  },
  {
    name: "containerMedium",
    value: "@container (min-width: 561px) and (max-width: 999px)",
  },
  {
    name: "containerLarge",
    value: "@container (min-width: 1000px)",
  },
] as const;

export const ResponsiveScale: Story = {
  render: () => (
    <FoundationPage
      title="Breakpoints"
      description="Breakpoint stories should show the actual responsive scale, not identical table rows."
    >
      <FoundationSection title="Viewport breakpoints">
        <styled.div display="flex" flexDirection="column" gap="3">
          {breakpointTokens.map((breakpoint) => (
            <styled.article
              key={breakpoint.name}
              backgroundColor="background.surface"
              borderColor="border.muted"
              borderRadius="panel"
              borderWidth="thin"
              display="grid"
              gap="4"
              gridTemplateColumns={{
                base: "1fr",
                md: "8rem minmax(0, 1fr) 9rem",
              }}
              padding="4"
            >
              <styled.div display="flex" flexDirection="column" gap="1">
                <TokenCode>{breakpoint.name}</TokenCode>
                <styled.span color="text.muted" fontSize="xs">
                  {breakpoint.usage}
                </styled.span>
              </styled.div>
              <styled.div
                backgroundColor="background.inset"
                borderRadius="pill"
                height="4"
                overflow="hidden"
                width="full"
              >
                <styled.div
                  backgroundColor="accent.default"
                  height="full"
                  style={{
                    width: `${(breakpoint.value / maxBreakpoint) * 100}%`,
                  }}
                />
              </styled.div>
              <styled.span color="text.default" fontFamily="mono" fontSize="sm">
                {breakpoint.value}px
              </styled.span>
            </styled.article>
          ))}
        </styled.div>
      </FoundationSection>

      <FoundationSection
        title="Container conditions"
        description="Container conditions are for components that adapt to their parent rather than the viewport."
      >
        <styled.div display="grid" gap="3">
          {containerConditions.map((condition) => (
            <styled.article
              key={condition.name}
              alignItems="center"
              backgroundColor="background.surface"
              borderColor="border.muted"
              borderRadius="panel"
              borderWidth="thin"
              display="grid"
              gap="3"
              gridTemplateColumns={{
                base: "1fr",
                md: "12rem minmax(0, 1fr)",
              }}
              padding="4"
            >
              <TokenCode>{condition.name}</TokenCode>
              <styled.span color="text.muted" fontFamily="mono" fontSize="sm">
                {condition.value}
              </styled.span>
            </styled.article>
          ))}
        </styled.div>
      </FoundationSection>
    </FoundationPage>
  ),
};
