import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Text } from "@/components/ui/text";
import { styled } from "@/styled-system/jsx";

import {
  FoundationPage,
  FoundationSection,
  TokenCode,
  tokenVar,
} from "./FoundationLayout";
import { fontFamilyTokens, fontSizeTokens } from "./tokenCatalog";

const meta = {
  title: "Foundations/Tokens/Typography",
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;
type TextVariant = "body" | "supporting" | "metadata";

const textStyles: Array<{
  name: TextVariant;
  token: string;
  sample: string;
}> = [
  {
    name: "body",
    token: "textStyles.body",
    sample: "There are 4 categories available to start discussions.",
  },
  {
    name: "supporting",
    token: "textStyles.supporting",
    sample: "Connect external MCP servers and add their tools to Robots.",
  },
  {
    name: "metadata",
    token: "textStyles.metadata",
    sample: "@southclaws · 18d · 0 replies",
  },
];

const fontWeights = [
  "normal",
  "medium",
  "semibold",
  "bold",
  "extrabold",
] as const;

export const TypeScale: Story = {
  render: () => (
    <FoundationPage
      title="Typography tokens"
      description="Typography foundations separate primitive font tokens from semantic UI prose. Text owns these three prose roles; headings and controls own their own metrics."
    >
      <FoundationSection title="Font families">
        <styled.div
          display="grid"
          gap="3"
          gridTemplateColumns="[repeat(auto-fit,minmax(13rem,1fr))]"
        >
          {fontFamilyTokens.map((font) => (
            <styled.article
              key={font.token}
              backgroundColor="background.surface"
              borderColor="border.muted"
              borderRadius="panel"
              borderWidth="thin"
              display="flex"
              flexDirection="column"
              gap="3"
              minW="0"
              padding="4"
            >
              <TokenCode>{font.token}</TokenCode>
              <styled.span
                color="text.default"
                fontSize="2xl"
                lineHeight="tight"
                style={{ fontFamily: tokenVar(font.token) }}
              >
                Storyden
              </styled.span>
              <styled.span color="text.muted" fontSize="sm">
                {font.usage}
              </styled.span>
            </styled.article>
          ))}
        </styled.div>
      </FoundationSection>

      <FoundationSection title="Font size scale">
        <styled.div display="flex" flexDirection="column" gap="2">
          {fontSizeTokens.map((size) => (
            <styled.article
              key={size}
              alignItems="baseline"
              backgroundColor="background.surface"
              borderColor="border.muted"
              borderRadius="panel"
              borderWidth="thin"
              display="grid"
              gap="4"
              gridTemplateColumns={{
                base: "1fr",
                md: "8rem minmax(0, 1fr)",
              }}
              padding="3"
            >
              <TokenCode>fontSizes.{size}</TokenCode>
              <styled.span
                color="text.default"
                lineHeight="tight"
                minW="0"
                style={{
                  fontSize: tokenVar(`fontSizes.${size}`),
                  overflowWrap: "anywhere",
                }}
              >
                The modern community surface
              </styled.span>
            </styled.article>
          ))}
        </styled.div>
      </FoundationSection>

      <FoundationSection title="Product text styles">
        <styled.div display="flex" flexDirection="column" gap="3">
          {textStyles.map((style) => (
            <styled.article
              key={style.name}
              backgroundColor="background.surface"
              borderColor="border.muted"
              borderRadius="panel"
              borderWidth="thin"
              display="grid"
              gap="3"
              gridTemplateColumns={{
                base: "1fr",
                md: "8rem minmax(0, 1fr)",
              }}
              padding="4"
            >
              <styled.div display="flex" flexDirection="column" gap="2">
                <TokenCode>{style.name}</TokenCode>
                <styled.span color="text.muted" fontSize="xs">
                  {style.token}
                </styled.span>
              </styled.div>
              <Text variant={style.name}>{style.sample}</Text>
            </styled.article>
          ))}
        </styled.div>
      </FoundationSection>

      <FoundationSection title="Weight">
        <styled.div
          display="grid"
          gap="3"
          gridTemplateColumns="[repeat(auto-fit,minmax(13rem,1fr))]"
        >
          {fontWeights.map((weight) => (
            <styled.article
              key={weight}
              backgroundColor="background.surface"
              borderColor="border.muted"
              borderRadius="panel"
              borderWidth="thin"
              display="flex"
              flexDirection="column"
              gap="2"
              padding="4"
            >
              <TokenCode>fontWeights.{weight}</TokenCode>
              <styled.span
                color="text.default"
                fontSize="lg"
                fontWeight={weight}
              >
                Storyden
              </styled.span>
            </styled.article>
          ))}
        </styled.div>
      </FoundationSection>
    </FoundationPage>
  ),
};
