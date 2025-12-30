import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import { TagNameList, TagReferenceList, Thread } from "@/api/openapi-schema";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Combotags, CombotagsItem } from "@/components/ui/combotags";

export type Props = {
  editing: boolean;
  initialTags?: TagReferenceList;
  onChange: (tags: TagNameList) => Promise<void>;
};

export function ThreadTagList(props: Props) {
  const currentTags = props.initialTags?.map((t) => t.name) ?? [];

  async function handleQuery(q: string): Promise<CombotagsItem[]> {
    const tags =
      (await handle(async () => {
        const { tags } = await tagList({ q });
        return tags.map((t) => t.name);
      })) ?? [];

    const filtered = tags.filter((t) => !currentTags.includes(t));

    return filtered.map((name) => ({ id: name, label: name }));
  }

  if (props.editing) {
    return (
      <>
        <Combotags
          initialValue={currentTags}
          onQuery={handleQuery}
          onChange={props.onChange}
        />
      </>
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

        return (
          <ThreadTagList
            editing={true}
            initialTags={initialTags}
            onChange={handleChange}
          />
        );
      }}
      control={control}
      name={name}
    />
  );
}
