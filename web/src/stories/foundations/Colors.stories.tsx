import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { styled } from "@/styled-system/jsx";
import { menu } from "@/styled-system/recipes";

import {
  FoundationPage,
  FoundationSection,
  TokenCode,
  compactTokenGrid,
  tokenVar,
} from "./FoundationLayout";
import {
  type LayerRoleTokenExample,
  type StatusTokenExample,
  type TokenExample,
  accentScaleTokens,
  borderTokens,
  layerRoleTokens,
  primitiveColorRamps,
  stateBackgroundTokens,
  statusTokens,
  textTokens,
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

const compactMenu = menu({ size: "xs" });

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
      backgroundColor="background.surface"
      borderColor="border.default"
      borderRadius="panel"
      borderWidth="thin"
      display="grid"
      gap="3"
      gridTemplateColumns="auto minmax(0, 1fr)"
      padding="3"
    >
      <styled.div
        alignItems="center"
        backgroundColor={preview === "content" ? "background.inset" : undefined}
        borderColor={preview === "border" ? undefined : "border.default"}
        borderRadius="control"
        borderWidth={preview === "border" ? "thick" : "thin"}
        color={preview === "content" ? undefined : "text.default"}
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
          <styled.strong color="text.default" fontSize="sm" fontWeight="medium">
            {name}
          </styled.strong>
          <TokenCode>{token}</TokenCode>
        </styled.div>
        <styled.p color="text.muted" fontSize="xs" lineHeight="relaxed">
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
      backgroundColor="background.surface"
      borderColor="border.default"
      borderRadius="panel"
      borderWidth="thin"
      display="flex"
      flexDirection="column"
      gap="3"
      padding="3"
    >
      <styled.div display="flex" flexDirection="column" gap="1">
        <styled.strong color="text.default" fontSize="sm" fontWeight="medium">
          {palette}
        </styled.strong>
        <TokenCode>colors.{palette}.1-12</TokenCode>
      </styled.div>
      <styled.div display="grid" gap="1" gridTemplateColumns="repeat(12, 1fr)">
        {steps.map((step) => (
          <styled.div
            key={step}
            borderColor="border.default"
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
      backgroundColor="background.surface"
      borderColor="border.default"
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
        <styled.strong color="text.default" fontSize="sm" fontWeight="medium">
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

function LayerRoleCard({
  name,
  usage,
  background,
  elevated,
}: LayerRoleTokenExample) {
  return (
    <styled.article
      borderColor="border.default"
      borderRadius="panel"
      borderStyle="solid"
      borderWidth="thin"
      display="flex"
      flexDirection="column"
      gap="2"
      padding="3"
      style={{
        backgroundColor: tokenVar(background),
        boxShadow: elevated ? tokenVar("shadows.overlay") : undefined,
      }}
    >
      <styled.div display="flex" flexDirection="column" gap="1">
        <styled.strong color="text.default" fontSize="sm" fontWeight="semibold">
          {name}
        </styled.strong>
        <styled.p color="text.muted" fontSize="sm" lineHeight="relaxed">
          {usage}
        </styled.p>
      </styled.div>
      <TokenCode>{background}</TokenCode>
    </styled.article>
  );
}

export const Layering: Story = {
  render: () => (
    <FoundationPage
      title="Layer roles"
      description="Canvas, Surface, Inset, and Overlay describe structural purpose. Accent and status colors remain separate from this hierarchy."
    >
      <FoundationSection
        title="Primary and secondary content"
        description="A Surface is independent. An Inset is subordinate and may sit inside a Surface or directly on the Canvas."
      >
        <styled.div
          backgroundColor="background.canvas"
          borderColor="border.default"
          borderRadius="panel"
          borderWidth="thin"
          color="text.default"
          display="grid"
          gap="4"
          gridTemplateColumns={{
            base: "1fr",
            md: "minmax(0, 2fr) minmax(15rem, 1fr)",
          }}
          padding="4"
        >
          <styled.article
            backgroundColor="background.surface"
            borderColor="border.default"
            borderRadius="lg"
            borderWidth="thin"
            color="text.default"
            display="flex"
            flexDirection="column"
            gap="2"
            padding="2"
          >
            <styled.div display="flex" flexDirection="column" gap="1">
              <styled.strong fontSize="md">Discussion report</styled.strong>
              <styled.span color="text.muted" fontSize="sm">
                An independent object on the application canvas.
              </styled.span>
            </styled.div>
            <styled.div
              backgroundColor="background.inset"
              borderColor="border.default"
              borderRadius="control"
              borderWidth="thin"
              color="text.default"
              display="flex"
              flexDirection="column"
              gap="1"
              padding="2"
            >
              <styled.strong fontSize="sm">
                Reported content preview
              </styled.strong>
              <styled.span color="text.muted" fontSize="xs">
                Embedded evidence uses Inset rather than another Surface.
              </styled.span>
            </styled.div>
          </styled.article>

          <styled.aside
            alignSelf="start"
            backgroundColor="background.inset"
            borderColor="border.default"
            borderRadius="lg"
            borderWidth="thin"
            color="text.default"
            display="flex"
            flexDirection="column"
            gap="1"
            padding="2"
          >
            <styled.strong fontSize="sm">Report summary</styled.strong>
            <styled.span color="text.muted" fontSize="xs">
              Supporting content can use Inset directly on Canvas.
            </styled.span>
          </styled.aside>
        </styled.div>
      </FoundationSection>

      <FoundationSection
        title="Temporary elevation"
        description="Overlay owns temporary content above normal flow. Shadow communicates elevation here, not ordinary content hierarchy."
      >
        <styled.div
          backgroundColor="background.canvas"
          borderColor="border.default"
          borderRadius="panel"
          borderWidth="thin"
          color="text.default"
          minHeight="64"
          overflow="hidden"
          padding="4"
          position="relative"
        >
          <styled.article
            backgroundColor="background.surface"
            borderColor="border.default"
            borderRadius="lg"
            borderWidth="thin"
            color="text.default"
            display="flex"
            flexDirection="column"
            gap="2"
            maxWidth="2xl"
            padding="2"
          >
            <styled.strong fontSize="md">Community settings</styled.strong>
            <styled.span color="text.muted" fontSize="sm">
              Normal document content remains below the temporary menu.
            </styled.span>
          </styled.article>

          <styled.div
            position="absolute"
            right="4"
            style={{ maxWidth: "calc(100% - 2rem)" }}
            top="16"
            width="40"
          >
            <styled.div
              backgroundColor="background.overlay"
              className={compactMenu.content}
              color="text.default"
              style={{ width: "min(10rem, calc(100vw - 4rem))" }}
            >
              <div className={compactMenu.itemGroup}>
                <strong className={compactMenu.itemGroupLabel}>Actions</strong>
                <styled.div
                  backgroundColor="background.controlHover"
                  className={compactMenu.item}
                >
                  Edit details
                </styled.div>
                <styled.div className={compactMenu.item} color="text.muted">
                  View activity
                </styled.div>
              </div>
            </styled.div>
          </styled.div>
        </styled.div>
      </FoundationSection>

      <FoundationSection
        title="Structure and meaning"
        description="Status colors communicate meaning without replacing the structural layer beneath the component."
      >
        <styled.div
          alignItems="start"
          backgroundColor="background.inset"
          borderColor="border.default"
          borderRadius="lg"
          borderWidth="thin"
          color="text.default"
          display="grid"
          gap="2"
          gridTemplateColumns="auto minmax(0, 1fr)"
          maxWidth="2xl"
          padding="2"
        >
          <styled.span
            backgroundColor="status.success.surface"
            borderColor="status.success.border"
            borderRadius="control"
            borderWidth="thin"
            color="status.success.content"
            fontSize="xs"
            fontWeight="semibold"
            paddingX="2"
            paddingY="1"
          >
            Success
          </styled.span>
          <styled.div display="flex" flexDirection="column" gap="1">
            <styled.strong fontSize="sm">Member import complete</styled.strong>
            <styled.span color="text.muted" fontSize="xs">
              Inset provides the structure; status tokens provide the meaning.
            </styled.span>
          </styled.div>
        </styled.div>
      </FoundationSection>
    </FoundationPage>
  ),
};

export const SemanticColors: Story = {
  render: () => (
    <FoundationPage
      title="Color tokens"
      description="Canonical color tokens describe structural layers first, then meaning and state. Product components should request semantic roles instead of palette steps."
    >
      <FoundationSection
        title="Structural backgrounds"
        description="Background tokens describe where a region sits in the document. Text, borders, and interaction states are selected independently by their own semantic role."
      >
        <div className={compactTokenGrid}>
          {layerRoleTokens.map((role) => (
            <LayerRoleCard key={role.name} {...role} />
          ))}
        </div>
      </FoundationSection>

      <FoundationSection
        title="Text emphasis"
        description="Text tokens describe emphasis rather than the layer containing the text."
      >
        <div className={compactTokenGrid}>
          {textTokens.map((token) => (
            <ColorTokenCard key={token.token} {...token} preview="content" />
          ))}
        </div>
      </FoundationSection>

      <FoundationSection
        title="Borders"
        description="Border strength is chosen independently from the background layer."
      >
        <div className={compactTokenGrid}>
          {borderTokens.map((token) => (
            <ColorTokenCard key={token.token} {...token} preview="border" />
          ))}
        </div>
      </FoundationSection>

      <FoundationSection
        title="Control and selection states"
        description="Interactive backgrounds are named for the state or control they belong to, not by visual intensity."
      >
        <div className={compactTokenGrid}>
          {stateBackgroundTokens.map((token) => (
            <ColorTokenCard key={token.token} {...token} />
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
