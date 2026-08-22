import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Text } from "@/components/ui/text";
import { LStack, styled } from "@/styled-system/jsx";

import {
  FoundationPage,
  FoundationSection,
  TokenCode,
} from "./FoundationLayout";

const meta = {
  title: "Foundations/Principles/Component Selection",
  parameters: {
    layout: "fullscreen",
    docs: {
      description: {
        component:
          "A semantic component chooser for Storyden product work. Start from the role the interface must communicate, then inspect the relevant component story and representative call sites before creating local visual code.",
      },
    },
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

type Selection = {
  need: string;
  startWith: string;
  use: string;
  avoid: string;
};

const structureSelections: Selection[] = [
  {
    need: "Page identity and actions",
    startWith: "PageHeader + PageHeading",
    use: "A screen title, return or hierarchy navigation, supporting context, and a small page-level action group.",
    avoid:
      "A feature-local header grid, manually aligned back button, or locally styled h1.",
  },
  {
    need: "Meaningful screen subdivision",
    startWith: "SectionHeading + Text",
    use: "An h2 section boundary with optional supporting copy on the page Canvas.",
    avoid:
      "A styled heading fragment, a Card solely to create separation, or a Surface title used as document hierarchy.",
  },
  {
    need: "Independent structured object",
    startWith: "Card / Surface",
    use: "An object with identity, description, metadata, media, controls, or repeated row and grid behavior.",
    avoid:
      "A custom bordered Box, a flattened domain object, or CardBox with many local style overrides.",
  },
  {
    need: "Compact arbitrary object",
    startWith: "CardBox",
    use: "A small independent object whose internal content does not need the structured Surface slots.",
    avoid:
      "Top-level page sections, settings groups, or a substitute for normal Canvas layout.",
  },
  {
    need: "Ordinary page or form grouping",
    startWith: "section + LStack / Grid",
    use: "Related fields or content that remain part of the page rather than becoming an independent object.",
    avoid:
      "Inventing a panel component, decorative wrapper, border, or background for every group.",
  },
];

const meaningSelections: Selection[] = [
  {
    need: "Lifecycle or execution state",
    startWith: "StatusBadge",
    use: "A domain-owned status label mapped to neutral, info, success, warning, or danger.",
    avoid:
      "Manually styling Badge with status token triples or using one Badge style for unrelated meanings.",
  },
  {
    need: "Category, tag, or compact label",
    startWith: "Badge",
    use: "Metadata that classifies or labels an object without representing its state.",
    avoid:
      "StatusBadge for decorative colour or a category that has no lifecycle meaning.",
  },
  {
    need: "Member identity",
    startWith: "MemberBadge",
    use: "Avatar plus the canonical member name or handle, with the established profile link behavior.",
    avoid:
      "Rebuilding identity from Avatar and raw profile text at each call site.",
  },
  {
    need: "Scan-friendly elapsed time",
    startWith: "RelativeTime",
    use: "Accessible approximate or strict relative time, paired with an exact domain timestamp when auditability matters.",
    avoid:
      "Direct date-fns formatting in product components or replacing a required exact timestamp.",
  },
];

const controlSelections: Selection[] = [
  {
    need: "Ordinary multiline input",
    startWith: "Textarea + FormControl",
    use: "Prose or instructions that need more space than Input, with the standard label and feedback composition.",
    avoid:
      "styled.textarea, local size recipes, or a specialist editor without a specialist editing model.",
  },
  {
    need: "Persistent contextual feedback",
    startWith: "Alert",
    use: "Risk, requirements, or consequences that remain visible while their condition applies.",
    avoid:
      "A locally coloured HStack or Box, field-level validation, or transient completion feedback.",
  },
  {
    need: "Task-level outcome",
    startWith: "Admonition",
    use: "A meaningful operation error, warning, or success state associated with the current task.",
    avoid:
      "A toast for critical recovery information or an Alert that cannot communicate acknowledgement.",
  },
  {
    need: "Empty collection or first-run state",
    startWith: "EmptyState",
    use: "A centred, concise explanation and a relevant next action when the product supports one.",
    avoid:
      "Loose placeholder copy, contribution language for operational tools, or arbitrary vertical centring code.",
  },
];

export const Guide: Story = {
  render: () => (
    <FoundationPage
      title="Choose components by meaning"
      description="The shared component is the default. Raw layout primitives arrange it; they do not replace its semantic role."
    >
      <FoundationSection
        title="Selection sequence"
        description="Complete this short pass before writing visual JSX for a product screen."
      >
        <styled.ol
          display="grid"
          gap="2"
          listStylePosition="inside"
          margin="0"
          maxW="3xl"
          padding="0"
        >
          <Text as="li">
            Name the semantic role the interface must express.
          </Text>
          <Text as="li">
            Find that role below, then read its component story and two product
            call sites.
          </Text>
          <Text as="li">
            Compose the screen from shared components and ordinary layout
            primitives.
          </Text>
          <Text as="li">
            Keep simple labels and arrangements inline. Extract a local
            component only when it owns a domain concept or substantial
            behavior.
          </Text>
          <Text as="li">
            If the same missing semantic role appears three times, stop and add
            it to the design system with its own decision-oriented story.
          </Text>
        </styled.ol>
      </FoundationSection>

      <SelectionSection title="Structure" selections={structureSelections} />
      <SelectionSection
        title="Identity and data"
        selections={meaningSelections}
      />
      <SelectionSection
        title="Controls and feedback"
        selections={controlSelections}
      />

      <FoundationSection
        title="What may remain feature-specific"
        description="A local component is valid when its boundary belongs to the product domain rather than a visual shortcut."
      >
        <LStack gap="3" maxW="3xl">
          <Text>
            Keep domain compositions such as a Trail overview, run card, or
            schedule field group local when they combine shared components with
            domain data, state, and behavior.
          </Text>
          <Text variant="supporting">
            Do not extract tiny count labels, status wrappers without domain
            mapping, or one-use layout fragments merely to shorten a parent
            file. Readability alone does not create a component boundary.
          </Text>
        </LStack>
      </FoundationSection>
    </FoundationPage>
  ),
};

function SelectionSection({
  title,
  selections,
}: {
  title: string;
  selections: Selection[];
}) {
  return (
    <FoundationSection title={title}>
      <styled.div>
        <styled.div
          display={{ base: "none", lg: "grid" }}
          gap="3"
          gridTemplateColumns="minmax(10rem, 0.8fr) minmax(12rem, 0.9fr) minmax(16rem, 1.4fr) minmax(16rem, 1.4fr)"
          paddingBlockEnd="2"
        >
          <Text variant="metadata" fontWeight="semibold">
            Need
          </Text>
          <Text variant="metadata" fontWeight="semibold">
            Start with
          </Text>
          <Text variant="metadata" fontWeight="semibold">
            Use when
          </Text>
          <Text variant="metadata" fontWeight="semibold">
            Avoid
          </Text>
        </styled.div>
        {selections.map((selection) => (
          <SelectionRow key={selection.need} {...selection} />
        ))}
      </styled.div>
    </FoundationSection>
  );
}

function SelectionRow({ need, startWith, use, avoid }: Selection) {
  return (
    <styled.article
      borderBlockStartColor="border.default"
      borderBlockStartWidth="thin"
      display="grid"
      gap="3"
      gridTemplateColumns={{
        base: "1fr",
        lg: "minmax(10rem, 0.8fr) minmax(12rem, 0.9fr) minmax(16rem, 1.4fr) minmax(16rem, 1.4fr)",
      }}
      paddingBlock="3"
    >
      <styled.strong fontSize="sm" fontWeight="semibold">
        {need}
      </styled.strong>
      <styled.div>
        <TokenCode>{startWith}</TokenCode>
      </styled.div>
      <Text variant="supporting">{use}</Text>
      <Text variant="metadata">Avoid: {avoid}</Text>
    </styled.article>
  );
}
