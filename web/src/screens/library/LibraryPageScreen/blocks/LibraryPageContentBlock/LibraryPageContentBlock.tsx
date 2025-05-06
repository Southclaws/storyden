import { useEffect, useState } from "react";
import { Controller } from "react-hook-form";

import { ContentComposer } from "src/components/content/ContentComposer/ContentComposer";

import { FormControl } from "@/components/ui/form/FormControl";

import { useLibraryPageContext } from "../../Context";
import { Form } from "../../form";
import { useAssets } from "../../useAssets";
import { useEditState } from "../../useEditState";

export function useLibraryPageContentBlock() {
  const { form } = useLibraryPageContext();

  const { content } = form.getValues();

  const [generatedContent, setGeneratedContent] = useState<string | undefined>(
    undefined,
  );

  function handleResetGeneratedContent() {
    setGeneratedContent(undefined);
  }

  useEffect(() => {
    if (content && generatedContent) {
      // if the content field changes, and it was previously using a controlled
      // value, reset this controlled value to move it back to uncontrolled.
      setGeneratedContent(undefined);
    }
  }, [form, content, generatedContent]);

  return {
    handleResetGeneratedContent,
    generatedContent,
  };
}

export function LibraryPageContentBlock() {
  const { form, node } = useLibraryPageContext();
  const { editing } = useEditState();
  const { handleAssetUpload } = useAssets(node);
  const { generatedContent } = useLibraryPageContentBlock();

  return (
    <FormControl>
      <Controller<Form>
        control={form.control}
        name="content"
        render={({ field }) => (
          <ContentComposer
            disabled={!editing}
            onAssetUpload={handleAssetUpload}
            initialValue={
              node.content ?? form.formState.defaultValues?.["content"]
            }
            value={generatedContent}
          />
        )}
      />
    </FormControl>
  );
}
