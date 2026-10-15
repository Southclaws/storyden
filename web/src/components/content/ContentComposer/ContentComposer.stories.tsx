import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { LStack } from "@/styled-system/jsx";

import { ContentComposer } from ".";

const meta = {
  title: "Compositions/Content/Composer",
  component: ContentComposer,
  parameters: {
    docs: {
      description: {
        component:
          "Use ContentComposer for persisted TipTap HTML that can switch between reading and editing. Reading uses the shared TipTap document schema without an editor instance. Editing loads the viewer's preferred rich-text or Markdown editor. Theme hooks: content-composer is the root, content-composer__content is read-only content, and content-composer__editor is the active rich-text or Markdown input. Use Text for plain labels and supporting copy; use ContentComposerField to bind editable content to a form.",
      },
    },
  },
  args: {
    disabled: true,
    initialValue:
      '<h2>Field notes</h2><p>Record <strong>what changed</strong> and leave a useful note for the next maker.</p><ul><li><p>Compare lighting and materials.</p></li><li><p>Keep a link to <a href="https://example.com/reference">the reference</a>.</p></li></ul><blockquote><p>Small details make an environment feel inhabited.</p></blockquote><pre><code>const experiment = "harbour";</code></pre>',
  },
} satisfies Meta<typeof ContentComposer>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Reading: Story = {};

export const Empty: Story = {
  args: { initialValue: "<p></p>" },
};

export const ReadAndEdit: Story = {
  render: (args) => {
    const [editing, setEditing] = useState(false);
    const [content, setContent] = useState(args.initialValue);
    return (
      <LStack gap="3">
        <Button onClick={() => setEditing(!editing)}>
          {editing ? "Read content" : "Edit content"}
        </Button>
        <ContentComposer
          {...args}
          disabled={!editing}
          initialValue={content}
          onChange={setContent}
        />
      </LStack>
    );
  },
};
