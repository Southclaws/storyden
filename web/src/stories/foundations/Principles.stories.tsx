import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { type ReactNode, useState } from "react";

import { Admonition } from "@/components/ui/admonition";
import * as Alert from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { Input } from "@/components/ui/input";
import { PageHeader } from "@/components/ui/page-header";
import * as Popover from "@/components/ui/popover";
import { SectionHeading } from "@/components/ui/section-heading";
import { Card } from "@/components/ui/surface";
import { Text } from "@/components/ui/text";
import { Grid, HStack, LStack, styled } from "@/styled-system/jsx";

import {
  FoundationPage,
  FoundationSection,
  TokenCode,
} from "./FoundationLayout";

const meta = {
  title: "Foundations/Principles/Storyden",
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        component:
          "The directional design guide for Storyden. Read these stories before choosing components or composing a new product screen.",
      },
    },
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

function Rule({ name, children }: { name: string; children: ReactNode }) {
  return (
    <styled.article
      borderBlockStartColor="border.default"
      borderBlockStartWidth="thin"
      display="flex"
      flexDirection="column"
      gap="1"
      paddingBlockStart="3"
    >
      <styled.strong fontSize="sm" fontWeight="semibold">
        {name}
      </styled.strong>
      <Text variant="supporting">{children}</Text>
    </styled.article>
  );
}

function DecisionRow({
  choice,
  use,
  avoid,
  example,
}: {
  choice: string;
  use: ReactNode;
  avoid: ReactNode;
  example: ReactNode;
}) {
  return (
    <styled.div
      alignItems={{ base: "start", lg: "center" }}
      borderBlockStartColor="border.default"
      borderBlockStartWidth="thin"
      display="grid"
      gap="3"
      gridTemplateColumns={{
        base: "1fr",
        lg: "8rem minmax(8rem, 0.8fr) minmax(12rem, 1fr) minmax(12rem, 1fr)",
      }}
      paddingBlock="3"
    >
      <styled.div>{example}</styled.div>
      <TokenCode>{choice}</TokenCode>
      <Text variant="supporting">{use}</Text>
      <Text variant="metadata">Avoid: {avoid}</Text>
    </styled.div>
  );
}

function Layer({
  name,
  token,
  children,
  background,
}: {
  name: string;
  token: string;
  children: ReactNode;
  background: "background.canvas" | "background.surface" | "background.inset";
}) {
  return (
    <styled.section
      backgroundColor={background}
      borderColor="border.default"
      borderRadius="panel"
      borderWidth={background === "background.inset" ? undefined : "thin"}
      display="flex"
      flexDirection="column"
      gap="3"
      padding="4"
    >
      <HStack justifyContent="space-between">
        <styled.strong fontSize="sm">{name}</styled.strong>
        <TokenCode>{token}</TokenCode>
      </HStack>
      {children}
    </styled.section>
  );
}

export const Overview: Story = {
  render: () => (
    <FoundationPage
      title="The Quiet Campfire"
      description="Storyden is compact, calm community software. The interface guides participation and curation without competing with the community's identity or content."
    >
      <FoundationSection
        title="Design priorities"
        description="When two valid solutions compete, decide in this order."
      >
        <Grid columns={{ base: 1, md: 2 }} gap="5">
          <Rule name="1. Meaning before appearance">
            Choose the semantic role first: action, navigation, content object,
            supporting region, status, or overlay. Styling follows that role.
          </Rule>
          <Rule name="2. Content before chrome">
            Community titles, discussion, member identity, and Library content
            should remain stronger than application controls.
          </Rule>
          <Rule name="3. Density without crowding">
            Use compact controls and close related groups. Preserve alignment,
            readable line height, and clear separation between unrelated work.
          </Rule>
          <Rule name="4. Neutral structure, selective meaning">
            Neutral layers express hierarchy. Accent and status colours only
            express interaction, selection, or meaning.
          </Rule>
          <Rule name="5. Composition before invention">
            Start with Storybook and existing UI components. A new visual
            component is a design-system decision, not a local convenience.
          </Rule>
          <Rule name="6. Accessibility is part of the visual language">
            Visible focus, readable contrast, keyboard operation, reduced
            motion, and stable responsive layouts are required states.
          </Rule>
        </Grid>
      </FoundationSection>

      <FoundationSection
        title="The implementation sequence"
        description="Use this sequence when building or reviewing a screen."
      >
        <styled.ol
          display="grid"
          gap="2"
          listStylePosition="inside"
          margin="0"
          padding="0"
        >
          <Text as="li">Identify the user's current task and return path.</Text>
          <Text as="li">
            Choose the page width and navigation model before styling content.
          </Text>
          <Text as="li">
            Establish page, section, and supporting-text hierarchy.
          </Text>
          <Text as="li">Assign canvas, surface, inset, and overlay roles.</Text>
          <Text as="li">
            Choose action emphasis based on consequence and workflow priority.
          </Text>
          <Text as="li">
            Verify empty, loading, disabled, error, mobile, and keyboard states.
          </Text>
        </styled.ol>
      </FoundationSection>
    </FoundationPage>
  ),
};

export const StructuralLayers: Story = {
  render: () => (
    <FoundationPage
      title="Structural layers"
      description="Layer tokens describe what a region is, not how dark, light, or elevated it happens to look in one theme."
    >
      <FoundationSection
        title="Canvas, Surface, Inset, Overlay"
        description="Storyden has two ordinary content layers and one temporary layer. Do not manufacture depth with arbitrary greys."
      >
        <Layer
          name="Canvas"
          token="background.canvas"
          background="background.canvas"
        >
          <Text variant="supporting">
            The application background. Page headings, settings forms, and
            top-level workflows may sit directly on it.
          </Text>
          <Layer
            name="Surface"
            token="background.surface"
            background="background.surface"
          >
            <Text variant="supporting">
              An independent object: a post, report, card, panel, or structured
              list item. Every Surface has a border; fill alone is not enough
              separation from the Canvas.
            </Text>
            <Layer
              name="Inset"
              token="background.inset"
              background="background.inset"
            >
              <Text variant="supporting">
                Supporting or embedded content: quoted posts, previews, metadata
                blocks, subcategory lists, code, and summaries.
              </Text>
            </Layer>
          </Layer>
        </Layer>

        <styled.div
          backgroundColor="background.canvas"
          borderColor="border.default"
          borderRadius="panel"
          borderWidth="thin"
          minHeight="48"
          padding="4"
          position="relative"
        >
          <Text variant="supporting">
            Normal document content continues here.
          </Text>
          <Popover.Root open>
            <Popover.Content
              borderRadius="overlay"
              display="flex"
              flexDirection="column"
              gap="2"
              maxW="xs"
              padding="3"
              position="absolute"
              right="4"
              top="4"
            >
              <HStack justifyContent="space-between">
                <styled.strong fontSize="sm">Overlay</styled.strong>
                <TokenCode>background.overlay</TokenCode>
              </HStack>
              <Text variant="supporting">
                Menus, popovers, dialogs, command palettes, and other temporary
                content above the document flow.
              </Text>
            </Popover.Content>
          </Popover.Root>
        </styled.div>
      </FoundationSection>

      <FoundationSection title="Layer rules">
        <Grid columns={{ base: 1, md: 2 }} gap="5">
          <Rule name="Surface means independent">
            Do not wrap every page section in a Surface. Settings screens and
            primary reading flows normally remain flush on the Canvas.
          </Rule>
          <Rule name="Inset means subordinate">
            An Inset can sit inside a Surface or directly on the Canvas, but do
            not nest Inset inside Inset.
          </Rule>
          <Rule name="Status is separate from structure">
            Alerts and Admonitions use Inset-like structure, then add danger,
            warning, success, or information tokens for meaning.
          </Rule>
          <Rule name="Shadows imply physical elevation">
            Routine hierarchy uses fill and border. Reserve overlay shadows for
            content that genuinely floats over the document.
          </Rule>
        </Grid>
      </FoundationSection>
    </FoundationPage>
  ),
};

export const ActionHierarchy: Story = {
  render: () => (
    <FoundationPage
      title="Action hierarchy"
      description="Button emphasis communicates priority within the current workflow. It is not decoration and it is not a global importance ranking."
    >
      <FoundationSection
        title="Choose a variant"
        description="Storyden defaults to compact subtle controls because most screens contain repeated, routine actions."
      >
        <DecisionRow
          example={<Button variant="subtle">Save</Button>}
          choice="subtle"
          use="The default for routine commands: save, create, connect, edit, share, and other reversible or familiar work."
          avoid="using it for every adjacent peer when one action must clearly lead."
        />
        <DecisionRow
          example={<Button variant="solid">Publish</Button>}
          choice="solid"
          use="The single decisive completion action in a focused flow: publish, confirm, authorise, finish setup."
          avoid="multiple solid buttons in one action group or using accent merely to add colour."
        />
        <DecisionRow
          example={<Button variant="outline">Preview</Button>}
          choice="outline"
          use="A visible secondary peer, alternate path, or persistent boundary needed against a variable background."
          avoid="using outline as the automatic default for every secondary action."
        />
        <DecisionRow
          example={<Button variant="ghost">Cancel</Button>}
          choice="ghost"
          use="Low-emphasis actions embedded in content or chrome: cancel, close, pagination, card tools, and navigation controls."
          avoid="important completion actions that need to remain discoverable at a glance."
        />
        <DecisionRow
          example={<Button variant="plain">More</Button>}
          choice="plain"
          use="Very dense rows and tree controls where the visible container would set the row height. The component retains a larger invisible hit target."
          avoid="ordinary buttons, forms, dialogs, or any action that needs a visible control boundary."
        />
      </FoundationSection>

      <FoundationSection
        title="Action intent"
        description="Variant expresses priority. Intent expresses meaning. Keep those decisions independent."
      >
        <HStack gap="2" flexWrap="wrap">
          <Button variant="solid">Continue</Button>
          <Button variant="subtle" intent="success">
            Approve
          </Button>
          <Button variant="outline" intent="destructive">
            Delete account
          </Button>
          <Button variant="solid" intent="destructive">
            Confirm deletion
          </Button>
        </HStack>
        <Text variant="supporting">
          Use a destructive outline or subtle action to enter a destructive
          flow. Use a solid destructive action only at the final irreversible
          confirmation. Success intent confirms an outcome; it should not
          replace the primary accent for ordinary progression.
        </Text>
      </FoundationSection>
    </FoundationPage>
  ),
};

export const TypographyAndScreens: Story = {
  render: () => (
    <FoundationPage
      title="UI typography and screen composition"
      description="UI typography is intentionally constrained. Rich community content uses the separate .typography system and its semantic heading tree."
    >
      <FoundationSection
        title="Product hierarchy"
        description="Heading level and visual role are related, but not interchangeable. Components own the correct semantic element for their context."
      >
        <styled.div
          borderBlockColor="border.default"
          borderBlockWidth="thin"
          display="flex"
          flexDirection="column"
          gap="6"
          paddingBlock="4"
        >
          <PageHeader
            title="System settings"
            description="Configure infrastructure, limits, and deployment behaviour."
            actions={<Button>Save</Button>}
          />
          <styled.section display="flex" flexDirection="column" gap="1">
            <SectionHeading>Client IP strategy</SectionHeading>
            <Text variant="supporting">
              Choose which proxy headers Storyden trusts when resolving a
              member's address.
            </Text>
          </styled.section>
          <Text>
            Body text explains the current task or value. Keep long UI copy
            within the readable measure and split genuinely separate concepts.
          </Text>
          <Text variant="metadata">Updated 12 minutes ago</Text>
        </styled.div>
      </FoundationSection>

      <FoundationSection title="Text roles">
        <DecisionRow
          example={<Text>Body text</Text>}
          choice="Text body"
          use="Primary UI prose, explanations, empty-state copy, and values that need normal reading emphasis."
          avoid="article or rich-text content; use the .typography system there."
        />
        <DecisionRow
          example={<Text variant="supporting">Supporting text</Text>}
          choice="Text supporting"
          use="Descriptions below page or section headings and secondary explanatory copy."
          avoid="timestamps, counts, IDs, or primary instructions."
        />
        <DecisionRow
          example={<Text variant="metadata">Metadata text</Text>}
          choice="Text metadata"
          use="Timestamps, counts, compact status details, identifiers, and low-priority annotations."
          avoid="essential instructions or long paragraphs."
        />
      </FoundationSection>

      <FoundationSection
        title="Page widths"
        description="The page chooses workspace width. Individual prose still owns its readable measure."
      >
        <Grid columns={{ base: 1, md: 3 }} gap="4">
          <Rule name="Content">
            The default for feeds, thread pages, forms, settings, and most
            reading or single-workflow screens.
          </Rule>
          <Rule name="Wide">
            Use for dashboards, Robots, complex tables, and tools that gain real
            utility from simultaneous columns.
          </Rule>
          <Rule name="Full">
            Reserve for canvas-like editors or data tools that explicitly own
            the viewport. Never let prose expand to the full width.
          </Rule>
        </Grid>
      </FoundationSection>

      <FoundationSection title="Representative composition">
        <Grid columns={{ base: 1, lg: 2 }} gap="4">
          <Card
            id="design-guidance-thread"
            title="Designing a calmer moderation queue"
            text="A discussion is an independent content object, so it uses the shared Surface card composition."
            url="/t/design-guidance-thread"
            disableAnchors
          >
            <Badge>Moderation</Badge>
          </Card>
          <CardBox as="section">
            <SectionHeading>Embedded summary</SectionHeading>
            <Text variant="supporting">
              CardBox is a compatibility shell for compact independent content.
              Prefer the Surface components when the content has a title,
              metadata, media, or repeated list/grid behaviour.
            </Text>
          </CardBox>
        </Grid>
      </FoundationSection>
    </FoundationPage>
  ),
};

export const ControlsAndFeedback: Story = {
  render: () => (
    <FoundationPage
      title="Controls, feedback, and status"
      description="Controls remain visually quiet until interaction. Feedback names the state and recovery without confusing hierarchy with meaning."
    >
      <FoundationSection title="Control surfaces">
        <Grid columns={{ base: 1, md: 2 }} gap="4">
          <styled.section display="flex" flexDirection="column" gap="2">
            <SectionHeading>Community name</SectionHeading>
            <Input defaultValue="Storyden" aria-label="Community name" />
            <Text variant="metadata">
              Inputs use background.control, not Surface. Inputs embedded in an
              Inset region may use the dedicated inset variant.
            </Text>
          </styled.section>
          <styled.section display="flex" flexDirection="column" gap="2">
            <SectionHeading>Status and categories</SectionHeading>
            <HStack gap="1" flexWrap="wrap">
              <Badge>General</Badge>
              <Badge colorPalette="green">Published</Badge>
              <Badge colorPalette="amber">Needs review</Badge>
              <Badge colorPalette="red">Blocked</Badge>
            </HStack>
            <Text variant="metadata">
              Badges label compact metadata. Palette must encode a real category
              or state, not provide decoration.
            </Text>
          </styled.section>
        </Grid>
      </FoundationSection>

      <FoundationSection title="Feedback hierarchy">
        <FeedbackExamples />
      </FoundationSection>
    </FoundationPage>
  ),
};

function FeedbackExamples() {
  const [showAdmonition, setShowAdmonition] = useState(true);

  return (
    <LStack gap="4">
      <LStack gap="3">
        <Alert.Root>
          <Alert.Icon asChild>
            <WarningIcon />
          </Alert.Icon>
          <Alert.Content>
            <Alert.Title>Access keys affect every API client</Alert.Title>
            <Alert.Description>
              Revoking this permission prevents members of the role from
              creating or using access keys.
            </Alert.Description>
          </Alert.Content>
        </Alert.Root>

        <Admonition
          value={showAdmonition}
          kind="failure"
          title="Changes could not be saved"
          onChange={setShowAdmonition}
        >
          Check the highlighted settings and try again.
        </Admonition>
      </LStack>

      <Grid columns={{ base: 1, md: 2 }} gap="5">
        <Rule name="Alert is persistent">
          Use Alert for contextual risk or guidance that must remain visible
          while its condition applies. It has no dismissal action.
        </Rule>
        <Rule name="Admonition is task feedback">
          Prefer Admonition for task-level errors, warnings, and meaningful
          success feedback. It stays with the relevant content and may be
          dismissed once acknowledged.
        </Rule>
      </Grid>

      <Text variant="supporting">
        Use inline field feedback for validation next to one control. Avoid
        toasts by default; reserve them for non-critical, transient confirmation
        with no action or recovery. Use a dialog only when focus must be
        protected or the action interrupted.
      </Text>
    </LStack>
  );
}
