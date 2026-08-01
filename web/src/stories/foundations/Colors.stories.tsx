import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { styled } from "@/styled-system/jsx";

import {
  FoundationPage,
  FoundationSection,
  TokenCode,
  compactTokenGrid,
  tokenVar,
} from "./FoundationLayout";
import {
  type StatusTokenExample,
  type TokenExample,
  accentScaleTokens,
  borderTokens,
  contentTokens,
  primitiveColorRamps,
  statusTokens,
  surfaceTokens,
  visibilityTokens,
} from "./tokenCatalog";

const meta = {
  title: "Foundations/Tokens/Colors",
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

function ColorTokenCard({
  name,
  token,
  usage,
  preview = "fill",
}: TokenExample & { preview?: "fill" | "content" | "border" }) {
  const variable = tokenVar(token);

  return (
    <styled.article
      alignItems="start"
      backgroundColor="surface.default"
      borderColor="border.subtle"
      borderRadius="panel"
      borderWidth="thin"
      display="grid"
      gap="3"
      gridTemplateColumns="auto minmax(0, 1fr)"
      padding="3"
    >
      <styled.div
        alignItems="center"
        backgroundColor={preview === "content" ? "surface.subtle" : undefined}
        borderColor={preview === "border" ? undefined : "border.default"}
        borderRadius="control"
        borderWidth={preview === "border" ? "thick" : "thin"}
        color={preview === "content" ? undefined : "content.default"}
        display="flex"
        fontSize="lg"
        fontWeight="semibold"
        height="12"
        justifyContent="center"
        style={{
          backgroundColor: preview === "fill" ? variable : undefined,
          borderColor: preview === "border" ? variable : undefined,
          color: preview === "content" ? variable : undefined,
        }}
        width="12"
      >
        {preview === "content" ? "Aa" : null}
      </styled.div>
      <styled.div display="flex" flexDirection="column" gap="2" minW="0">
        <styled.div display="flex" flexDirection="column" gap="1">
          <styled.strong
            color="content.default"
            fontSize="sm"
            fontWeight="medium"
          >
            {name}
          </styled.strong>
          <TokenCode>{token}</TokenCode>
        </styled.div>
        <styled.p color="content.subtle" fontSize="xs" lineHeight="relaxed">
          {usage}
        </styled.p>
      </styled.div>
    </styled.article>
  );
}

function PrimitiveRamp({ palette }: { palette: string }) {
  const steps = Array.from({ length: 12 }, (_, index) => index + 1);

  return (
    <styled.article
      backgroundColor="surface.default"
      borderColor="border.subtle"
      borderRadius="panel"
      borderWidth="thin"
      display="flex"
      flexDirection="column"
      gap="3"
      padding="3"
    >
      <styled.div display="flex" flexDirection="column" gap="1">
        <styled.strong
          color="content.default"
          fontSize="sm"
          fontWeight="medium"
        >
          {palette}
        </styled.strong>
        <TokenCode>colors.{palette}.1-12</TokenCode>
      </styled.div>
      <styled.div display="grid" gap="1" gridTemplateColumns="repeat(12, 1fr)">
        {steps.map((step) => (
          <styled.div
            key={step}
            borderColor="border.subtle"
            borderRadius="xs"
            borderWidth="thin"
            height="12"
            style={{ backgroundColor: tokenVar(`colors.${palette}.${step}`) }}
            title={`colors.${palette}.${step}`}
          />
        ))}
      </styled.div>
    </styled.article>
  );
}

function StatusCard({ name, surface, content, border }: StatusTokenExample) {
  return (
    <styled.article
      alignItems="center"
      backgroundColor="surface.default"
      borderColor="border.subtle"
      borderRadius="panel"
      borderWidth="thin"
      display="grid"
      gap="3"
      gridTemplateColumns={{
        base: "1fr",
        md: "minmax(8rem, 1fr) minmax(0, 2fr)",
      }}
      padding="3"
    >
      <styled.div display="flex" flexDirection="column" gap="1">
        <styled.strong
          color="content.default"
          fontSize="sm"
          fontWeight="medium"
        >
          {name}
        </styled.strong>
        <TokenCode>{surface}</TokenCode>
      </styled.div>
      <styled.div
        borderRadius="control"
        borderStyle="solid"
        borderWidth="thin"
        fontSize="sm"
        fontWeight="medium"
        paddingX="3"
        paddingY="2"
        style={{
          backgroundColor: tokenVar(surface),
          borderColor: tokenVar(border),
          color: tokenVar(content),
        }}
      >
        {name} state
      </styled.div>
    </styled.article>
  );
}

export const SemanticColors: Story = {
  render: () => (
    <FoundationPage
      title="Color tokens"
      description="Canonical color tokens describe jobs in the interface: canvas, surface, content, border, interaction, status, accent, and domain visibility."
    >
      <FoundationSection
        title="Surfaces"
        description="Use surface tokens for canvas-adjacent UI layers instead of choosing raw neutral palette steps."
      >
        <div className={compactTokenGrid}>
          {surfaceTokens.map((token) => (
            <ColorTokenCard key={token.token} {...token} />
          ))}
        </div>
      </FoundationSection>

      <FoundationSection title="Content">
        <div className={compactTokenGrid}>
          {contentTokens.map((token) => (
            <ColorTokenCard key={token.token} {...token} preview="content" />
          ))}
        </div>
      </FoundationSection>

      <FoundationSection title="Borders">
        <div className={compactTokenGrid}>
          {borderTokens.map((token) => (
            <ColorTokenCard key={token.token} {...token} preview="border" />
          ))}
        </div>
      </FoundationSection>

      <FoundationSection title="Status">
        <styled.div display="grid" gap="3">
          {statusTokens.map((token) => (
            <StatusCard key={token.name} {...token} />
          ))}
        </styled.div>
      </FoundationSection>

      <FoundationSection title="Accent">
        <div className={compactTokenGrid}>
          {accentScaleTokens.map((token) => (
            <ColorTokenCard key={token.token} {...token} />
          ))}
        </div>
      </FoundationSection>

      <FoundationSection title="Visibility">
        <styled.div display="grid" gap="3">
          {visibilityTokens.map((token) => (
            <StatusCard key={token.name} {...token} />
          ))}
        </styled.div>
      </FoundationSection>

      <FoundationSection
        title="Primitive ramps"
        description="Primitive palette ramps are available for token authors, but product components should prefer semantic tokens."
      >
        <styled.div display="grid" gap="3">
          {primitiveColorRamps.map((palette) => (
            <PrimitiveRamp key={palette} palette={palette} />
          ))}
        </styled.div>
      </FoundationSection>
    </FoundationPage>
  ),
};
