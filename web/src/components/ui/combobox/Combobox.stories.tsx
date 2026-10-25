import { Portal, createListCollection } from "@ark-ui/react";
import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useId } from "react";
import { useForm } from "react-hook-form";
import { expect, userEvent, waitFor, within } from "storybook/test";

import { Box, LStack } from "@/styled-system/jsx";

import { Button } from "../button";
import { FormControl } from "../form-control";
import { FormLabel } from "../form-label";
import { CheckIcon } from "../icons/Check";
import { SelectIcon } from "../icons/Select";
import { Input } from "../input";
import { Menu } from "../menu/Menu";
import * as Popover from "../popover";
import { Popover as ClosedPopover } from "../popover/Popover";
import { Textarea } from "../textarea";

import * as Combobox from ".";
import { ComboboxField } from "./Combobox.field";

const items = [
  { label: "General", value: "general" },
  { label: "Announcements", value: "announcements" },
  { label: "Support", value: "support" },
  { label: "Showcase", value: "showcase" },
];

type ComboboxStoryArgs = {
  size?: "sm" | "md" | "lg";
};

const meta = {
  title: "Components/Forms/Combobox",
  parameters: {
    docs: {
      description: {
        component:
          "A searchable selection control for longer or dynamic collections. Use Select for a short known option set and Input for free-form values. The trigger, input, and icon inherit one shared size. Portal the Positioner when used inside an overlay; the shared content recipe keeps its option list above the owning menu or popover without a local z-index override.",
      },
    },
  },
  argTypes: {
    size: {
      control: "select",
      options: ["sm", "md", "lg"],
    },
  },
} satisfies Meta<ComboboxStoryArgs>;

export default meta;

type Story = StoryObj<ComboboxStoryArgs>;

function Example({
  size,
  options = items,
}: ComboboxStoryArgs & { options?: typeof items }) {
  const collection = createListCollection({ items: options });
  return (
    <Combobox.Root
      collection={collection}
      maxW="sm"
      positioning={{ sameWidth: true }}
      size={size}
    >
      <Combobox.Label>Category</Combobox.Label>
      <Combobox.Control>
        <Combobox.Input placeholder="Select a category" asChild>
          <Input size={size} />
        </Combobox.Input>
        <Combobox.Trigger asChild>
          <button type="button" aria-label="Open">
            <SelectIcon />
          </button>
        </Combobox.Trigger>
      </Combobox.Control>
      <Portal>
        <Combobox.Positioner>
          <Combobox.Content>
            <Combobox.List>
              <Combobox.ItemGroup>
                {options.map((item) => (
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
      </Portal>
    </Combobox.Root>
  );
}

export const Playground: Story = {
  render: (args) => <Example {...args} />,
};

export const Sizes: Story = {
  render: () => (
    <LStack gap="4">
      {(["sm", "md", "lg"] as const).map((size) => (
        <Example key={size} size={size} />
      ))}
    </LStack>
  ),
};

export const LongLabels: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "Options stay on one line and truncate with an ellipsis when space is limited. The full label remains searchable and available to assistive technology, with room for the selected indicator at every size.",
      },
    },
  },
  render: () => (
    <LStack gap="4" w="64" alignItems="stretch">
      {(["sm", "md", "lg"] as const).map((size) => (
        <Example
          key={size}
          size={size}
          options={[
            { label: "General", value: "general" },
            {
              label:
                "Outline test — environment art notebook and field references",
              value: "notebook",
            },
            {
              label:
                "AnUnbrokenPageNameThatMustNotOverflowTheOptionOrHideItsSelectedIndicator",
              value: "unbroken",
            },
          ]}
        />
      ))}
    </LStack>
  ),
};

export const InMenu: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "A portalled picker inside a menu uses the shared combobox elevation. Its options remain above the menu even where their edges overlap, without a consumer z-index override.",
      },
    },
  },
  render: (args) => (
    <Menu
      defaultOpen
      trigger={<Button variant="outline">Category settings</Button>}
      contentProps={{ w: "80" }}
    >
      <Box p="4">
        <Example {...args} />
      </Box>
    </Menu>
  ),
  play: async ({ canvasElement }) => {
    const body = within(canvasElement.ownerDocument.body);
    await userEvent.click(await body.findByRole("button", { name: "Open" }));
    const popup = await body.findByRole("listbox");
    await waitFor(() => {
      const rect = popup.getBoundingClientRect();
      for (const [x, y] of [
        [rect.left + 4, rect.top + 4],
        [rect.right - 4, rect.top + 4],
        [rect.left + 4, rect.bottom - 4],
        [rect.right - 4, rect.bottom - 4],
      ] as const) {
        expect(
          popup.contains(canvasElement.ownerDocument.elementFromPoint(x, y)),
        ).toBe(true);
      }
    });
    await userEvent.click(body.getByRole("option", { name: "Support" }));
    await expect(body.getByRole("combobox")).toHaveValue("Support");
    await waitFor(() => expect(popup).not.toBeVisible());
  },
};

export const InPopover: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "The same portalled picker also layers above a popover. Supporting content below the input deliberately overlaps the options to exercise their stacking and pointer interaction.",
      },
    },
  },
  render: (args) => (
    <ClosedPopover
      defaultOpen
      trigger={<Button variant="outline">Category settings</Button>}
      contentProps={{ w: "80", p: "4" }}
    >
      <Popover.Title>Category settings</Popover.Title>
      <Example {...args} />
      <Popover.Description>
        Choose the category for new posts.
      </Popover.Description>
    </ClosedPopover>
  ),
  play: InMenu.play,
};

export const FollowingField: Story = {
  parameters: {
    docs: {
      description: {
        story:
          "Select a category, then Tab or click into the instructions. Closing the option list keeps focus on the following field once you have moved there.",
      },
    },
  },
  render: function FollowingField() {
    const form = useForm({ defaultValues: { category: "", instruction: "" } });
    const instructionID = useId();
    return (
      <LStack gap="4" maxW="sm" alignItems="stretch">
        <FormControl>
          <FormLabel>Category</FormLabel>
          <ComboboxField
            control={form.control}
            name="category"
            items={items}
            placeholder="Select a category"
            ariaLabel="Category"
          />
        </FormControl>
        <FormControl>
          <FormLabel id={instructionID}>Instructions</FormLabel>
          <Textarea
            aria-labelledby={instructionID}
            {...form.register("instruction")}
          />
        </FormControl>
      </LStack>
    );
  },
};
