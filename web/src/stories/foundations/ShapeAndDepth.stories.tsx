import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { styled } from "@/styled-system/jsx";

import {
  FoundationPage,
  FoundationSection,
  TokenCode,
  compactTokenGrid,
  tokenValue,
  tokenVar,
} from "./FoundationLayout";
import { blurTokens, radiusTokens, shadowTokens } from "./tokenCatalog";

const meta = {
  title: "Foundations/Tokens/Shape and Depth",
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

export const ShapeAndDepth: Story = {
  render: () => (
    <FoundationPage
      title="Shape and depth tokens"
      description="Shape and depth tokens should be restrained and semantic enough that components do not invent local corner, elevation, or blur values."
    >
      <FoundationSection title="Radii">
        <div className={compactTokenGrid}>
          {radiusTokens.map((token) => (
            <styled.article
              key={token}
              alignItems="center"
              backgroundColor="background.surface"
              borderColor="border.muted"
              borderRadius="panel"
              borderWidth="thin"
              display="grid"
              gap="3"
              gridTemplateColumns="auto minmax(0, 1fr)"
              padding="3"
            >
              <styled.div
                backgroundColor="accent.subtle"
                borderColor="accent.default"
                borderStyle="solid"
                borderWidth="medium"
                height="16"
                style={{ borderRadius: tokenVar(token) }}
                width="16"
              />
              <styled.div
                display="flex"
                flexDirection="column"
                gap="1"
                minW="0"
              >
                <TokenCode>{token}</TokenCode>
                <styled.span
                  color="text.muted"
                  fontFamily="mono"
                  fontSize="xs"
                  style={{ overflowWrap: "anywhere" }}
                >
                  {tokenValue(token)}
                </styled.span>
              </styled.div>
            </styled.article>
          ))}
        </div>
      </FoundationSection>

      <FoundationSection title="Shadows">
        <div className={compactTokenGrid}>
          {shadowTokens.map((token) => (
            <styled.article
              key={token}
              alignItems="center"
              backgroundColor="background.inset"
              borderColor="border.muted"
              borderRadius="panel"
              borderWidth="thin"
              display="grid"
              gap="3"
              gridTemplateColumns="auto minmax(0, 1fr)"
              padding="5"
            >
              <styled.div
                backgroundColor="background.overlay"
                borderColor="border.muted"
                borderRadius="panel"
                borderWidth="thin"
                height="16"
                style={{ boxShadow: tokenVar(token) }}
                width="16"
              />
              <styled.div
                display="flex"
                flexDirection="column"
                gap="1"
                minW="0"
              >
                <TokenCode>{token}</TokenCode>
                <styled.span
                  color="text.muted"
                  fontFamily="mono"
                  fontSize="xs"
                  style={{ overflowWrap: "anywhere" }}
                >
                  {tokenValue(token)}
                </styled.span>
              </styled.div>
            </styled.article>
          ))}
        </div>
      </FoundationSection>

      <FoundationSection title="Blurs">
        <div className={compactTokenGrid}>
          {blurTokens.map((token) => (
            <styled.article
              key={token}
              backgroundColor="background.surface"
              borderColor="border.muted"
              borderRadius="panel"
              borderWidth="thin"
              display="flex"
              flexDirection="column"
              gap="3"
              padding="3"
            >
              <styled.div
                borderColor="border.muted"
                borderRadius="overlay"
                borderWidth="thin"
                height="20"
                overflow="hidden"
                position="relative"
                style={{
                  backgroundImage: `repeating-linear-gradient(90deg, ${tokenVar("colors.status.success.surface")} 0 0.75rem, ${tokenVar("colors.status.warning.surface")} 0.75rem 1.5rem, ${tokenVar("colors.status.danger.surface")} 1.5rem 2.25rem)`,
                }}
              >
                <styled.div
                  borderLeftColor="border.default"
                  borderLeftWidth="thin"
                  bottom="0"
                  position="absolute"
                  right="0"
                  style={{
                    backdropFilter: `blur(${tokenVar(token)})`,
                    backgroundColor:
                      "color-mix(in srgb, var(--colors-background-surface) 20%, transparent)",
                    WebkitBackdropFilter: `blur(${tokenVar(token)})`,
                  }}
                  top="0"
                  width="60%"
                />
              </styled.div>
              <styled.div
                display="flex"
                flexDirection="column"
                gap="1"
                minW="0"
              >
                <TokenCode>{token}</TokenCode>
                <styled.span
                  color="text.muted"
                  fontFamily="mono"
                  fontSize="xs"
                  style={{ overflowWrap: "anywhere" }}
                >
                  {tokenValue(token)}
                </styled.span>
              </styled.div>
            </styled.article>
          ))}
        </div>
      </FoundationSection>
    </FoundationPage>
  ),
};
