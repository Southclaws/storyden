import { CategorySelect } from "@/components/category/CategorySelect/CategorySelect";
import { ContentFormField } from "@/components/content/ContentComposer/ContentField";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { Input } from "@/components/ui/input";
import { HStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Props, useCreateThreadFromResult } from "./useCreateThreadFromResult";

export function CreateThreadFromResultScreen(props: Props) {
  const { form, contentHTML, handlers } = useCreateThreadFromResult(props);

  return (
    <form className={lstack()} onSubmit={handlers.handleSubmit}>
      <FormControl>
        <FormHelperText>Thread title</FormHelperText>
        <Input placeholder="Title" {...form.register("title")} />
        <FormErrorText>{form.formState.errors.title?.message}</FormErrorText>
      </FormControl>

      <FormControl>
        <FormHelperText>Category</FormHelperText>
        <CategorySelect name="category" control={form.control} />
        <FormErrorText>{form.formState.errors.category?.message}</FormErrorText>
      </FormControl>

      <FormControl>
        <FormHelperText>Thread content</FormHelperText>
        <ContentFormField
          name="content"
          control={form.control}
          initialValue={contentHTML}
        />
        <FormErrorText>{form.formState.errors.content?.message}</FormErrorText>
      </FormControl>

      <WStack>
        <Button
          size="sm"
          variant="outline"
          type="button"
          onClick={props.onFinish}
        >
          Cancel
        </Button>

        <Button size="sm" variant="solid">
          Post
        </Button>
      </WStack>
    </form>
  );
}
