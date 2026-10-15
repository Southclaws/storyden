import { useState } from "react";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import { TagNameList } from "@/api/openapi-schema";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/multi-select-picker";

import type { Props } from "./ThreadTagList";

export function ThreadTagListEditing(props: Omit<Props, "editing">) {
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
      disabled={props.disabled}
    />
  );
}
