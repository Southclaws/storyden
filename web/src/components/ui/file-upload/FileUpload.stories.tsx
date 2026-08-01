import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { Button } from "../button";
import { DeleteIcon } from "../icons/Delete";

import * as FileUpload from ".";

const meta = {
  title: "UI/File Upload",
  component: FileUpload.Root,
} satisfies Meta<typeof FileUpload.Root>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Dropzone: Story = {
  render: () => (
    <FileUpload.Root maxFiles={3} accept="image/*" maxW="lg">
      <FileUpload.Label>Upload assets</FileUpload.Label>
      <FileUpload.Dropzone>
        <FileUpload.Trigger asChild>
          <Button variant="outline">Choose files</Button>
        </FileUpload.Trigger>
        <span>Drag image files here</span>
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
                  <Button size="xs" variant="ghost">
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
  ),
};
