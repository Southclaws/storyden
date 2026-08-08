import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { styled } from "@/styled-system/jsx";

import {
  FoundationPage,
  FoundationSection,
  TokenCode,
  tokenValue,
  tokenVar,
} from "./FoundationLayout";
import {
  type TokenExample,
  layoutSizeTokens,
  spacingTokens,
} from "./tokenCatalog";

const meta = {
  title: "Foundations/Tokens/Spacing",
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

function SpacingRow({ name }: { name: string }) {
  const token = `spacing.${name}`;
  const variable = tokenVar(token);

  return (
    <styled.article
      alignItems="center"
      backgroundColor="background.surface"
      borderColor="border.muted"
      borderRadius="panel"
      borderWidth="thin"
      display="grid"
      gap="4"
      gridTemplateColumns={{
        base: "1fr",
        md: "8rem minmax(0, 30rem) 7rem",
      }}
      padding="3"
    >
      <TokenCode>{token}</TokenCode>
      <styled.div
        backgroundColor="background.inset"
        borderRadius="control"
        height="5"
        overflow="hidden"
        width="full"
      >
        <styled.div
          backgroundColor="text.default"
          height="full"
          minW={name === "0" ? "0" : "0.5"}
          style={{ width: variable }}
        />
      </styled.div>
      <styled.span color="text.muted" fontFamily="mono" fontSize="xs">
        {tokenValue(token)}
      </styled.span>
    </styled.article>
  );
}

function LayoutSizeRow({ name, token, usage }: TokenExample) {
  const variable = tokenVar(token);

  return (
    <styled.article
      backgroundColor="background.surface"
      borderColor="border.muted"
      borderRadius="panel"
      borderWidth="thin"
      display="flex"
      flexDirection="column"
      gap="2"
      padding="3"
    >
      <styled.div
        alignItems={{ base: "start", md: "center" }}
        display="grid"
        gap="2"
        gridTemplateColumns={{ base: "1fr", md: "12rem minmax(0, 1fr) auto" }}
      >
        <styled.div display="flex" flexDirection="column" gap="1">
          <styled.strong color="text.default" fontSize="sm" fontWeight="medium">
            {name}
          </styled.strong>
          <TokenCode>{token}</TokenCode>
        </styled.div>
        <styled.p color="text.muted" fontSize="xs" lineHeight="relaxed">
          {usage}
        </styled.p>
        <styled.span color="text.subtle" fontFamily="mono" fontSize="xs">
          {tokenValue(token)}
        </styled.span>
      </styled.div>
      <styled.div
        backgroundColor="background.inset"
        borderRadius="control"
        height="6"
        overflow="hidden"
        width="full"
      >
        <styled.div
          backgroundColor="text.default"
          height="full"
          maxW="full"
          style={{ width: variable }}
        />
      </styled.div>
    </styled.article>
  );
}

export const SpacingScale: Story = {
  render: () => (
    <FoundationPage
      title="Spacing and layout tokens"
      description="Primitive spacing tokens are for local rhythm. Semantic layout size tokens are for repeated product dimensions."
    >
      <FoundationSection
        title="Spacing scale"
        description="Bars are rendered at token width against a fixed track so increments are visually comparable."
      >
        <styled.div display="flex" flexDirection="column" gap="2" width="full">
          {spacingTokens.map((name) => (
            <SpacingRow key={name} name={name} />
          ))}
        </styled.div>
      </FoundationSection>

      <FoundationSection
        title="Layout sizes"
        description="Use these for recurring application dimensions instead of choosing raw spacing or size steps each time. Bars use one shared track so their widths remain comparable."
      >
        <styled.div display="flex" flexDirection="column" gap="2" width="full">
          {layoutSizeTokens.map((token) => (
            <LayoutSizeRow key={token.token} {...token} />
          ))}
        </styled.div>
      </FoundationSection>
    </FoundationPage>
  ),
};
