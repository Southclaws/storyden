import { createListCollection } from "@ark-ui/react";
import { parseColor } from "@ark-ui/react/color-picker";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { type CSSProperties, type ReactNode, useState } from "react";

import { Button } from "@/components/ui/button";
import * as Checkbox from "@/components/ui/checkbox";
import * as ColorPicker from "@/components/ui/color-picker";
import * as Combobox from "@/components/ui/combobox";
import { DatePicker } from "@/components/ui/date-picker";
import * as FileUpload from "@/components/ui/file-upload";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormFeedback } from "@/components/ui/form-feedback";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { ColourPipetteIcon } from "@/components/ui/icons/Colour";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { LayoutGridIcon } from "@/components/ui/icons/LayoutGrid";
import { LayoutListIcon } from "@/components/ui/icons/LayoutList";
import { LayoutTableIcon } from "@/components/ui/icons/LayoutTable";
import { SearchIcon } from "@/components/ui/icons/Search";
import { SelectIcon } from "@/components/ui/icons/Select";
import { Input } from "@/components/ui/input";
import { InputGroup } from "@/components/ui/input-group";
import {
  MultiSelectPicker,
  type MultiSelectPickerItem,
} from "@/components/ui/multi-select-picker";
import { NumberInput } from "@/components/ui/number-input";
import { PinInput } from "@/components/ui/pin-input";
import * as RadioGroup from "@/components/ui/radio-group";
import * as Select from "@/components/ui/select";
import { Slider } from "@/components/ui/slider";
import { Switch } from "@/components/ui/switch";
import * as ToggleGroup from "@/components/ui/toggle-group";
import { styled } from "@/styled-system/jsx";

const sizes = ["sm", "md", "lg"] as const;
const buttonVariants = ["subtle", "outline", "ghost", "solid"] as const;
const inputVariants = ["outline", "ghost"] as const;
const selectVariants = ["outline", "ghost"] as const;
const toggleVariants = ["outline", "ghost"] as const;

type ControlSize = (typeof sizes)[number];

const categories = [
  { label: "Discussion", value: "discussion" },
  { label: "Library page", value: "library" },
  { label: "Collection", value: "collection" },
  { label: "Profile", value: "profile" },
];

const categoryCollection = createListCollection({ items: categories });

const tags = [
  { label: "Design", value: "design", colour: "#6d5dfc" },
  { label: "Frontend", value: "frontend", colour: "#27b981" },
  { label: "Accessibility", value: "accessibility" },
  { label: "Moderation", value: "moderation" },
  { label: "Performance", value: "performance" },
  { label: "Storybook", value: "storybook" },
];

const selectedTags = tags.slice(0, 2);
const colorPresets = [
  "hsl(10, 81%, 59%)",
  "hsl(60, 81%, 59%)",
  "hsl(100, 81%, 59%)",
  "hsl(175, 81%, 59%)",
  "hsl(220, 81%, 59%)",
  "hsl(280, 81%, 59%)",
  "hsl(350, 81%, 59%)",
];

const meta = {
  title: "Compositions/Forms/Component Audit",
  parameters: {
    layout: "fullscreen",
  },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

function AuditShell({
  children,
  description,
  title,
}: {
  children: ReactNode;
  description: string;
  title: string;
}) {
  return (
    <styled.main
      boxSizing="border-box"
      color="content.default"
      display="flex"
      flexDirection="column"
      gap="8"
      marginX="auto"
      maxW="7xl"
      padding={{ base: "4", md: "8" }}
      width="full"
    >
      <styled.header display="flex" flexDirection="column" gap="2" maxW="3xl">
        <styled.h1 fontSize="2xl" fontWeight="semibold" lineHeight="tight">
          {title}
        </styled.h1>
        <styled.p color="content.subtle" fontSize="md" lineHeight="relaxed">
          {description}
        </styled.p>
      </styled.header>
      {children}
    </styled.main>
  );
}

function Section({
  children,
  description,
  title,
}: {
  children: ReactNode;
  description?: string;
  title: string;
}) {
  return (
    <styled.section display="flex" flexDirection="column" gap="4">
      <styled.div display="flex" flexDirection="column" gap="1">
        <styled.h2 fontSize="xl" fontWeight="semibold">
          {title}
        </styled.h2>
        {description && (
          <styled.p color="content.subtle" fontSize="sm">
            {description}
          </styled.p>
        )}
      </styled.div>
      {children}
    </styled.section>
  );
}

function Matrix({
  children,
  columns,
}: {
  children: ReactNode;
  columns: string;
}) {
  const columnStyle = { "--audit-columns": columns } as CSSProperties;

  return (
    <styled.div
      borderTopColor="border.subtle"
      borderTopWidth="thin"
      display="grid"
      gap="0"
      width="full"
    >
      <styled.div
        alignItems="end"
        borderBottomColor="border.subtle"
        borderBottomWidth="thin"
        color="content.subtle"
        display="grid"
        fontSize="sm"
        fontWeight="medium"
        gap="3"
        gridTemplateColumns={{ base: "1fr", lg: "var(--audit-columns)" }}
        paddingY="2"
        style={columnStyle}
      >
        <styled.div />
        {children}
      </styled.div>
    </styled.div>
  );
}

function MatrixRow({
  children,
  columns,
  label,
  meta,
}: {
  children: ReactNode;
  columns: string;
  label: string;
  meta?: string;
}) {
  const columnStyle = { "--audit-columns": columns } as CSSProperties;

  return (
    <styled.div
      alignItems="center"
      borderBottomColor="border.subtle"
      borderBottomWidth="thin"
      display="grid"
      gap="3"
      gridTemplateColumns={{ base: "1fr", lg: "var(--audit-columns)" }}
      paddingY="4"
      style={columnStyle}
    >
      <styled.div display="flex" flexDirection="column" gap="1">
        <styled.strong fontSize="sm" fontWeight="medium">
          {label}
        </styled.strong>
        {meta && (
          <styled.code color="content.muted" fontFamily="mono" fontSize="xs">
            {meta}
          </styled.code>
        )}
      </styled.div>
      {children}
    </styled.div>
  );
}

function FitCell({ children }: { children: ReactNode }) {
  return <styled.div justifySelf="start">{children}</styled.div>;
}

function ExampleSelect({
  size,
  variant = "outline",
}: {
  size: ControlSize;
  variant?: (typeof selectVariants)[number];
}) {
  return (
    <Select.Root collection={categoryCollection} size={size} variant={variant}>
      <Select.Control>
        <Select.Trigger>
          <Select.ValueText placeholder="Choose type" />
          <SelectIcon />
        </Select.Trigger>
      </Select.Control>
      <Select.Positioner>
        <Select.Content>
          {categories.map((item) => (
            <Select.Item key={item.value} item={item}>
              <Select.ItemText>{item.label}</Select.ItemText>
              <Select.ItemIndicator>
                <CheckIcon />
              </Select.ItemIndicator>
            </Select.Item>
          ))}
        </Select.Content>
      </Select.Positioner>
      <Select.HiddenSelect />
    </Select.Root>
  );
}

function ExampleCombobox({ size }: { size: ControlSize }) {
  return (
    <Combobox.Root collection={categoryCollection} size={size}>
      <Combobox.Control>
        <Combobox.Input placeholder="Search categories" asChild>
          <Input size={size} />
        </Combobox.Input>
        <Combobox.Trigger asChild>
          <button type="button" aria-label="Open category list">
            <SelectIcon />
          </button>
        </Combobox.Trigger>
      </Combobox.Control>
      <Combobox.Positioner>
        <Combobox.Content>
          <Combobox.List>
            <Combobox.ItemGroup>
              {categories.map((item) => (
                <Combobox.Item key={item.value} item={item}>
                  <Combobox.ItemText>{item.label}</Combobox.ItemText>
                  <Combobox.ItemIndicator>
                    <CheckIcon />
                  </Combobox.ItemIndicator>
                </Combobox.Item>
              ))}
            </Combobox.ItemGroup>
          </Combobox.List>
        </Combobox.Content>
      </Combobox.Positioner>
    </Combobox.Root>
  );
}

function ExampleMultiSelect({ size }: { size: ControlSize }) {
  const [value, setValue] = useState<MultiSelectPickerItem[]>(selectedTags);
  const [results, setResults] = useState<MultiSelectPickerItem[]>(tags);

  return (
    <MultiSelectPicker
      allowNewValues
      autoColour
      inputPlaceholder="Select tags"
      queryResults={results}
      size={size}
      value={value}
      onChange={async (next) => setValue(next)}
      onQuery={(query) =>
        setResults(
          tags.filter((item) =>
            item.label.toLowerCase().includes(query.toLowerCase()),
          ),
        )
      }
    />
  );
}

function ExampleRadioGroup({ size }: { size: ControlSize }) {
  return (
    <RadioGroup.Root defaultValue="public" orientation="horizontal" size={size}>
      {["public", "members", "private"].map((value) => (
        <RadioGroup.Item key={value} value={value}>
          <RadioGroup.ItemControl />
          <RadioGroup.ItemText>
            {value.charAt(0).toUpperCase() + value.slice(1)}
          </RadioGroup.ItemText>
          <RadioGroup.ItemHiddenInput />
        </RadioGroup.Item>
      ))}
    </RadioGroup.Root>
  );
}

function ExampleToggleGroup({
  size,
  variant = "outline",
}: {
  size: ControlSize;
  variant?: (typeof toggleVariants)[number];
}) {
  return (
    <ToggleGroup.Root defaultValue={["list"]} size={size} variant={variant}>
      <ToggleGroup.Item value="list" aria-label="List">
        <LayoutListIcon />
      </ToggleGroup.Item>
      <ToggleGroup.Item value="grid" aria-label="Grid">
        <LayoutGridIcon />
      </ToggleGroup.Item>
      <ToggleGroup.Item value="table" aria-label="Table">
        <LayoutTableIcon />
      </ToggleGroup.Item>
    </ToggleGroup.Root>
  );
}

function ExampleColorPicker() {
  return (
    <ColorPicker.Root defaultValue={parseColor("hsl(175, 81%, 39%)")}>
      <ColorPicker.Context>
        {(api) => (
          <>
            <ColorPicker.Control>
              <ColorPicker.ChannelInput channel="hex" asChild>
                <Input maxW="44" />
              </ColorPicker.ChannelInput>
              <ColorPicker.Trigger asChild>
                <IconButton variant="outline" aria-label="Open color picker">
                  <ColorPicker.Swatch value={api.value} />
                </IconButton>
              </ColorPicker.Trigger>
            </ColorPicker.Control>
            <ColorPicker.Positioner>
              <ColorPicker.Content>
                <styled.div display="flex" flexDirection="column" gap="3">
                  <ColorPicker.Area>
                    <ColorPicker.AreaBackground />
                    <ColorPicker.AreaThumb />
                  </ColorPicker.Area>
                  <styled.div alignItems="center" display="flex" gap="3">
                    <ColorPicker.EyeDropperTrigger asChild>
                      <IconButton
                        aria-label="Pick a color"
                        size="sm"
                        variant="outline"
                      >
                        <ColourPipetteIcon />
                      </IconButton>
                    </ColorPicker.EyeDropperTrigger>
                    <styled.div display="flex" flex="1" flexDirection="column">
                      <ColorPicker.ChannelSlider channel="hue">
                        <ColorPicker.ChannelSliderTrack />
                        <ColorPicker.ChannelSliderThumb />
                      </ColorPicker.ChannelSlider>
                    </styled.div>
                  </styled.div>
                  <ColorPicker.SwatchGroup>
                    {colorPresets.map((color) => (
                      <ColorPicker.SwatchTrigger key={color} value={color}>
                        <ColorPicker.Swatch value={color} />
                      </ColorPicker.SwatchTrigger>
                    ))}
                  </ColorPicker.SwatchGroup>
                </styled.div>
              </ColorPicker.Content>
            </ColorPicker.Positioner>
          </>
        )}
      </ColorPicker.Context>
      <ColorPicker.HiddenInput />
    </ColorPicker.Root>
  );
}

function ExampleFileUpload() {
  return (
    <FileUpload.Root accept="image/*" maxFiles={3}>
      <FileUpload.Label>Upload assets</FileUpload.Label>
      <FileUpload.Dropzone>
        <FileUpload.Trigger asChild>
          <Button variant="outline">Choose files</Button>
        </FileUpload.Trigger>
        <styled.span color="fg.muted" fontSize="sm">
          Drag image files here
        </styled.span>
      </FileUpload.Dropzone>
      <FileUpload.ItemGroup>
        <FileUpload.Context>
          {(api) =>
            api.acceptedFiles.map((file) => (
              <FileUpload.Item key={file.name} file={file}>
                <FileUpload.ItemPreview />
                <FileUpload.ItemName />
                <FileUpload.ItemSizeText />
                <FileUpload.ItemDeleteTrigger asChild>
                  <Button size="sm" variant="ghost">
                    <DeleteIcon />
                    Remove
                  </Button>
                </FileUpload.ItemDeleteTrigger>
              </FileUpload.Item>
            ))
          }
        </FileUpload.Context>
      </FileUpload.ItemGroup>
      <FileUpload.HiddenInput />
    </FileUpload.Root>
  );
}

export const TextEntrySizing: Story = {
  render: () => {
    const coreColumns =
      "5rem minmax(9rem, 1fr) minmax(12rem, 1fr) minmax(12rem, 1fr) minmax(12rem, 1fr) minmax(12rem, 1fr)";
    const lookupColumns =
      "5rem minmax(14rem, 1fr) minmax(14rem, 1fr) minmax(18rem, 1fr)";

    return (
      <AuditShell
        title="Form text-entry sizing"
        description="Compare controls that commonly appear in the same form row. The row height should step consistently across sm, md, and lg."
      >
        <Section
          title="Core controls"
          description="Button, Input, InputGroup, Select, and NumberInput are most likely to sit in the same horizontal form row."
        >
          <Matrix columns={coreColumns}>
            <styled.div>Button</styled.div>
            <styled.div>Input</styled.div>
            <styled.div>Input group</styled.div>
            <styled.div>Select</styled.div>
            <styled.div>Number</styled.div>
          </Matrix>
          {sizes.map((size) => (
            <MatrixRow
              key={size}
              columns={coreColumns}
              label={size.toUpperCase()}
              meta={`size="${size}"`}
            >
              <Button size={size} variant="subtle">
                Save
              </Button>
              <Input placeholder="storyden" size={size} />
              <InputGroup
                size={size}
                startElement={<SearchIcon />}
                endElement={<SelectIcon />}
              >
                <Input placeholder="Search" size={size} />
              </InputGroup>
              <ExampleSelect size={size} />
              <NumberInput
                defaultValue="12"
                inputProps={{ "aria-label": `Number input ${size}` }}
                size={size}
              />
            </MatrixRow>
          ))}
        </Section>

        <Section
          title="Lookup and code controls"
          description="Combobox, MultiSelectPicker, and PinInput have different content density, so they get a separate sizing row."
        >
          <Matrix columns={lookupColumns}>
            <styled.div>Combobox</styled.div>
            <styled.div>Multi select</styled.div>
            <styled.div>Pin input</styled.div>
          </Matrix>
          {sizes.map((size) => (
            <MatrixRow
              key={size}
              columns={lookupColumns}
              label={size.toUpperCase()}
              meta={`size="${size}"`}
            >
              <ExampleCombobox size={size} />
              <ExampleMultiSelect size={size} />
              <PinInput length={4} size={size}>
                Code
              </PinInput>
            </MatrixRow>
          ))}
        </Section>
      </AuditShell>
    );
  },
};

export const ControlVariants: Story = {
  render: () => (
    <AuditShell
      title="Form control variants"
      description="Compare variants for the common button, input, select, number input, and toggle surfaces."
    >
      <Section title="Buttons">
        <styled.div
          display="grid"
          gap="3"
          gridTemplateColumns={{
            base: "1fr",
            md: "repeat(4, minmax(8rem, 1fr))",
          }}
        >
          {buttonVariants.map((variant) => (
            <Button key={variant} variant={variant}>
              {variant}
            </Button>
          ))}
        </styled.div>
      </Section>
      <Section title="Text controls">
        <styled.div
          display="grid"
          gap="3"
          gridTemplateColumns={{ base: "1fr", md: "repeat(3, 1fr)" }}
        >
          {inputVariants.map((variant) => (
            <Input
              key={variant}
              placeholder={`Input ${variant}`}
              variant={variant}
            />
          ))}
          {selectVariants.map((variant) => (
            <ExampleSelect key={variant} size="sm" variant={variant} />
          ))}
          <NumberInput
            defaultValue="12"
            inputProps={{ "aria-label": "Outline number input" }}
            variant="outline"
          />
          <NumberInput
            defaultValue="24"
            inputProps={{ "aria-label": "Ghost number input" }}
            variant="ghost"
          />
        </styled.div>
      </Section>
      <Section title="Toggle groups">
        <styled.div display="flex" flexWrap="wrap" gap="3">
          {toggleVariants.map((variant) => (
            <ExampleToggleGroup key={variant} size="sm" variant={variant} />
          ))}
        </styled.div>
      </Section>
    </AuditShell>
  ),
};

export const SelectionControls: Story = {
  render: () => {
    const booleanColumns =
      "5rem minmax(12rem, 1fr) minmax(12rem, 1fr) minmax(18rem, 2fr)";
    const choiceColumns = "5rem minmax(24rem, 2fr) minmax(12rem, 1fr)";

    return (
      <AuditShell
        title="Form selection controls"
        description="Compare the controls that represent choice, boolean state, segmented state, and bounded numeric ranges."
      >
        <Section title="Boolean and range controls">
          <Matrix columns={booleanColumns}>
            <styled.div>Checkbox</styled.div>
            <styled.div>Switch</styled.div>
            <styled.div>Slider</styled.div>
          </Matrix>
          {sizes.map((size) => (
            <MatrixRow
              key={size}
              columns={booleanColumns}
              label={size.toUpperCase()}
              meta={`size="${size}"`}
            >
              <Checkbox.Checkbox defaultChecked size={size}>
                Published
              </Checkbox.Checkbox>
              <Switch defaultChecked size={size}>
                Notify
              </Switch>
              <Slider
                defaultValue={[64]}
                marks={[
                  { value: 0, label: "0" },
                  { value: 50, label: "50" },
                  { value: 100, label: "100" },
                ]}
                max={100}
                min={0}
                size={size}
              />
            </MatrixRow>
          ))}
        </Section>

        <Section title="Choice and segmented controls">
          <Matrix columns={choiceColumns}>
            <styled.div>Radio group</styled.div>
            <styled.div>Toggle group</styled.div>
          </Matrix>
          {sizes.map((size) => (
            <MatrixRow
              key={size}
              columns={choiceColumns}
              label={size.toUpperCase()}
              meta={`size="${size}"`}
            >
              <ExampleRadioGroup size={size} />
              <FitCell>
                <ExampleToggleGroup size={size} />
              </FitCell>
            </MatrixRow>
          ))}
        </Section>
      </AuditShell>
    );
  },
};

export const FieldStatesAndComplexInputs: Story = {
  render: () => (
    <AuditShell
      title="Form fields and complex inputs"
      description="Compare labels, helper text, error text, disabled/invalid states, and the larger picker-style form components."
    >
      <Section title="Field states">
        <styled.div
          display="grid"
          gap="4"
          gridTemplateColumns={{ base: "1fr", md: "repeat(3, 1fr)" }}
        >
          <FormControl>
            <FormLabel>Community name</FormLabel>
            <Input defaultValue="Storyden" />
            <FormHelperText>
              Shown in navigation and browser title.
            </FormHelperText>
          </FormControl>
          <FormControl>
            <FormLabel>Slug</FormLabel>
            <Input aria-invalid defaultValue="storyden!!" />
            <FormFeedback error="Use lowercase letters, numbers, and hyphens." />
          </FormControl>
          <FormControl>
            <FormLabel>Invite code</FormLabel>
            <Input disabled defaultValue="private-beta" />
            <FormErrorText>
              Disabled field styles should stay legible.
            </FormErrorText>
          </FormControl>
        </styled.div>
      </Section>
      <Section title="Complex controls">
        <styled.div
          alignItems="start"
          display="grid"
          gap="5"
          gridTemplateColumns={{ base: "1fr", lg: "repeat(2, minmax(0, 1fr))" }}
        >
          <FormControl>
            <FormLabel>Publish date</FormLabel>
            <DatePicker />
            <FormHelperText>
              Closed date-picker trigger beside its input.
            </FormHelperText>
          </FormControl>
          <FormControl>
            <FormLabel>Theme colour</FormLabel>
            <ExampleColorPicker />
            <FormHelperText>
              Closed colour-picker trigger beside hex input.
            </FormHelperText>
          </FormControl>
          <FormControl>
            <FormLabel>Tags</FormLabel>
            <ExampleMultiSelect size="sm" />
            <FormHelperText>
              Badge density inside a button-shaped trigger.
            </FormHelperText>
          </FormControl>
          <ExampleFileUpload />
        </styled.div>
      </Section>
    </AuditShell>
  ),
};
