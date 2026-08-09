import { createListCollection } from "@ark-ui/react";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useForm } from "react-hook-form";

import { Button } from "@/components/ui/button";
import { CheckboxField } from "@/components/ui/checkbox";
import { ColorPickerField } from "@/components/ui/color-picker";
import { ComboboxField } from "@/components/ui/combobox";
import { DatePickerField } from "@/components/ui/date-picker";
import { FormControl } from "@/components/ui/form-control";
import { FormLabel } from "@/components/ui/form-label";
import { HeadingInputField } from "@/components/ui/heading-input";
import { NumberInputField } from "@/components/ui/number-input";
import { PinInputField } from "@/components/ui/pin-input";
import { RadioGroupField } from "@/components/ui/radio-group";
import { SelectField } from "@/components/ui/select";
import { SliderField } from "@/components/ui/slider";
import { ToggleGroupField } from "@/components/ui/toggle-group";
import { styled } from "@/styled-system/jsx";

const sectionCollection = createListCollection({
  items: [
    { label: "Discussion", value: "discussion" },
    { label: "Library", value: "library" },
  ],
});

type FormValues = {
  title: string;
  enabled: boolean;
  visibility: string;
  section?: string;
  model?: string;
  count?: number;
  volume?: number;
  code: string;
  date?: string;
  colour?: string;
  kinds?: string[];
};

const meta = {
  title: "Compositions/Forms/React Hook Form Fields",
  parameters: { layout: "fullscreen" },
} satisfies Meta;

export default meta;

type Story = StoryObj<typeof meta>;

export const FieldAdapters: Story = {
  render: function FieldAdapterStory() {
    const form = useForm<FormValues>({
      defaultValues: {
        title: "Community settings",
        enabled: true,
        visibility: "public",
        section: "discussion",
        model: "openai/gpt-5",
        count: 12,
        volume: 50,
        code: "1234",
        date: "2026-08-05T00:00:00.000Z",
        colour: "#31a89f",
        kinds: ["thread"],
      },
    });

    return (
      <styled.form
        display="grid"
        gap="5"
        gridTemplateColumns={{ base: "1fr", md: "repeat(2, minmax(0, 1fr))" }}
        marginX="auto"
        maxW="5xl"
        padding="8"
        width="full"
        onSubmit={form.handleSubmit(() => undefined)}
      >
        <FormControl>
          <FormLabel>Page title</FormLabel>
          <HeadingInputField
            control={form.control}
            name="title"
            aria-label="Page title"
          />
        </FormControl>

        <FormControl>
          <FormLabel>Section</FormLabel>
          <SelectField
            control={form.control}
            name="section"
            collection={sectionCollection}
            placeholder="Choose section"
          />
        </FormControl>

        <FormControl>
          <FormLabel>Model</FormLabel>
          <ComboboxField
            control={form.control}
            name="model"
            items={[
              { label: "GPT-5", value: "openai/gpt-5" },
              { label: "Claude", value: "anthropic/claude" },
            ]}
            placeholder="Choose model"
            ariaLabel="Choose model"
          />
        </FormControl>

        <FormControl>
          <FormLabel>Visibility</FormLabel>
          <RadioGroupField
            control={form.control}
            name="visibility"
            items={[
              { label: "Public", value: "public" },
              { label: "Private", value: "private" },
            ]}
          />
        </FormControl>

        <FormControl>
          <FormLabel>Count</FormLabel>
          <NumberInputField
            control={form.control}
            name="count"
            ariaLabel="Count"
            min={0}
          />
        </FormControl>

        <FormControl>
          <SliderField control={form.control} name="volume" label="Volume" />
        </FormControl>

        <FormControl>
          <FormLabel>Verification code</FormLabel>
          <PinInputField control={form.control} name="code" length={4} />
        </FormControl>

        <FormControl>
          <FormLabel>Expiry date</FormLabel>
          <DatePickerField control={form.control} name="date" />
        </FormControl>

        <FormControl>
          <FormLabel>Accent colour</FormLabel>
          <ColorPickerField
            control={form.control}
            name="colour"
            ariaLabel="Accent colour"
          />
        </FormControl>

        <FormControl>
          <FormLabel>Content kinds</FormLabel>
          <ToggleGroupField
            control={form.control}
            name="kinds"
            items={[
              { label: "Threads", value: "thread", description: "Threads" },
              { label: "Pages", value: "page", description: "Pages" },
            ]}
          />
        </FormControl>

        <CheckboxField control={form.control} name="enabled">
          Enable notifications
        </CheckboxField>

        <styled.div display="flex" gap="2" justifyContent="end">
          <Button type="button" variant="ghost" onClick={() => form.reset()}>
            Reset
          </Button>
          <Button type="submit">Save</Button>
        </styled.div>
      </styled.form>
    );
  },
};
