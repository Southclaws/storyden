import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { useState } from "react";

import { ThemeEditorPanel } from "./ThemeEditorPanel";
import {
  type ThemeEditorDocument,
  themeDocumentsSignature,
} from "./theme-editor-model";

const initialDocuments: ThemeEditorDocument[] = [
  {
    key: "css-1",
    kind: "stylesheet",
    label: "CSS 1",
    source:
      ":root {\n  --sd-color-accent: rebeccapurple;\n  --sd-radius-panel: 6px;\n}\n",
    savedSource:
      ":root {\n  --sd-color-accent: rebeccapurple;\n  --sd-radius-panel: 6px;\n}\n",
  },
  {
    key: "js-1",
    kind: "script",
    label: "JS 1",
    source: "document.documentElement.dataset.communityTheme = 'active';\n",
    savedSource:
      "document.documentElement.dataset.communityTheme = 'active';\n",
  },
];

const meta = {
  title: "Compositions/Admin/Theme Editor",
  component: ThemeEditorPanel,
  parameters: { layout: "fullscreen" },
} satisfies Meta<typeof ThemeEditorPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const LiveEditor: Story = {
  args: {
    documents: initialDocuments,
    selectedKey: "css-1",
    dirty: false,
    busy: false,
    runtimeDisabled: false,
    onSelect: () => undefined,
    onAdd: () => undefined,
    onChange: () => undefined,
    onMove: () => undefined,
    onDelete: () => undefined,
    onSave: () => undefined,
    onExit: () => undefined,
  },
  render: () => <InteractiveEditor />,
};

export const Empty: Story = {
  args: {
    documents: [],
    dirty: false,
    busy: false,
    runtimeDisabled: false,
    onSelect: () => undefined,
    onAdd: () => undefined,
    onChange: () => undefined,
    onMove: () => undefined,
    onDelete: () => undefined,
    onSave: () => undefined,
    onExit: () => undefined,
  },
};

export const ValidationError: Story = {
  args: {
    ...Empty.args,
    documents: initialDocuments,
    selectedKey: "css-1",
    dirty: true,
    error: "CSS 1 exceeds the 1 MiB asset limit.",
  },
};

function InteractiveEditor() {
  const [documents, setDocuments] = useState(initialDocuments);
  const [selectedKey, setSelectedKey] = useState("css-1");
  const [baseline] = useState(themeDocumentsSignature(initialDocuments));
  return (
    <ThemeEditorPanel
      documents={documents}
      selectedKey={selectedKey}
      dirty={themeDocumentsSignature(documents) !== baseline}
      busy={false}
      runtimeDisabled={false}
      onSelect={setSelectedKey}
      onAdd={() => undefined}
      onChange={(key, source) =>
        setDocuments((current) =>
          current.map((item) =>
            item.key === key ? { ...item, source } : item,
          ),
        )
      }
      onMove={() => undefined}
      onDelete={(key) =>
        setDocuments((current) => current.filter((item) => item.key !== key))
      }
      onSave={() => undefined}
      onExit={() => undefined}
    />
  );
}
