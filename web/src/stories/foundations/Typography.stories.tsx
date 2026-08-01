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
type TextSize = "xs" | "sm" | "md" | "lg" | "xl" | "2xl";

const textStyles: Array<{
  name: TextSize;
  token: string;
  sample: string;
}> = [
  {
    name: "xs",
    token: "textStyle xs",
    sample: "@southclaws • 18d • 0 replies",
  },
  {
    name: "sm",
    token: "textStyle sm",
    sample: "Discussion categories",
  },
  {
    name: "md",
    token: "textStyle md",
    sample: "There are 4 categories available to start discussions.",
  },
  {
    name: "lg",
    token: "textStyle lg",
    sample: "Check out the new Storyden command line interface!",
  },
  {
    name: "xl",
    token: "textStyle xl",
    sample: "Related discussions",
  },
  {
    name: "2xl",
    token: "textStyle 2xl",
    sample: "Create discussion",
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
      description="Typography foundations separate primitive font tokens from product text styles. Components should usually consume text styles."
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
              backgroundColor="surface.default"
              borderColor="border.subtle"
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
                color="content.default"
                fontSize="2xl"
                lineHeight="tight"
                style={{ fontFamily: tokenVar(font.token) }}
              >
                Storyden
              </styled.span>
              <styled.span color="content.subtle" fontSize="sm">
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
              backgroundColor="surface.default"
              borderColor="border.subtle"
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
                color="content.default"
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
              backgroundColor="surface.default"
              borderColor="border.subtle"
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
                <styled.span color="content.subtle" fontSize="xs">
                  {style.token}
                </styled.span>
              </styled.div>
              <Text color="content.default" size={style.name}>
                {style.sample}
              </Text>
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
              backgroundColor="surface.default"
              borderColor="border.subtle"
              borderRadius="panel"
              borderWidth="thin"
              display="flex"
              flexDirection="column"
              gap="2"
              padding="4"
            >
              <TokenCode>fontWeights.{weight}</TokenCode>
              <styled.span
                color="content.default"
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
