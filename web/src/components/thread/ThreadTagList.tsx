import { useState } from "react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import { TagNameList, TagReferenceList } from "@/api/openapi-schema";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/MultiSelectPicker";

export type Props = {
  editing: boolean;
  initialTags?: TagReferenceList;
  onChange: (tags: TagNameList) => Promise<void>;
};

export function ThreadTagList(props: Props) {
  const [queryResults, setQueryResults] = useState<MultiSelectPickerItem[]>(
    [],
  );

  const currentTags: MultiSelectPickerItem[] =
    props.initialTags?.map((t) => ({
      label: t.name,
      value: t.name,
    })) ?? [];

  function handleQuery(q: string) {
    handle(async () => {
      const { tags } = await tagList({ q });
      const filtered = tags.filter(
        (t) => !currentTags.some((ct) => ct.value === t.name),
      );
      setQueryResults(
        filtered.map((t) => ({
          label: t.name,
          value: t.name,
        })),
      );
    });
  }

  async function handleChange(items: MultiSelectPickerItem[]) {
    const tagNames = items.map((item) => item.value);
    await props.onChange(tagNames);
  }

  if (props.editing) {
    return (
      <MultiSelectPicker
        value={currentTags}
        onChange={handleChange}
        onQuery={handleQuery}
        queryResults={queryResults}
        allowNewValues={true}
        inputPlaceholder="Add tags..."
        autoColour={true}
        size="sm"
      />
    );
  }

  if (props.initialTags?.length === 0) {
    return null;
  }

  return <TagBadgeList tags={props.initialTags ?? []} />;
}

type TagListFieldProps<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  initialTags?: TagReferenceList;
};

export function TagListField<T extends FieldValues>({
  control,
  name,
  initialTags,
}: TagListFieldProps<T>) {
  return (
    <Controller<T>
      render={({ field }) => {
        async function handleChange(tags: string[]) {
          field.onChange(tags);
        }

        const fieldTags =
          field.value?.map((name: string) => ({ name })) || initialTags || [];

        return (
          <ThreadTagList
            editing={true}
            initialTags={fieldTags}
            onChange={handleChange}
          />
        );
      }}
      control={control}
      name={name}
    />
  );
}
