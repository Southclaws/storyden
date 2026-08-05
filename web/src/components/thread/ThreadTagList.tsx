import { useState } from "react";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import { TagNameList, TagReferenceList } from "@/api/openapi-schema";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/multi-select-picker";

export type Props = {
  editing: boolean;
  disabled?: boolean;
  initialTags?: TagReferenceList;
  onChange: (tags: TagNameList) => Promise<void>;
};

export function ThreadTagList(props: Props) {
  const [queryResults, setQueryResults] = useState<MultiSelectPickerItem[]>([]);

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
        triggerProps={{ disabled: props.disabled }}
      />
    );
  }

  if (props.initialTags?.length === 0) {
    return null;
  }

  return <TagBadgeList tags={props.initialTags ?? []} />;
}
