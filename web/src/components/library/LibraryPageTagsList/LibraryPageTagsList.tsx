import { useRef, useState } from "react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import { InstanceCapability, Node, TagNameList } from "@/api/openapi-schema";
import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Combotags, CombotagsHandle } from "@/components/ui/combotags";
import { useLibraryMutation } from "@/lib/library/library";
import { useCapability } from "@/lib/settings/capabilities";
import { HStack } from "@/styled-system/jsx";

export type Props<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  editing: boolean;
  node: Node;
};

export function useLibraryPageTagsList<T extends FieldValues>(props: Props<T>) {
  const isSuggestEnabled = useCapability(InstanceCapability.gen_ai);

  const { suggestTags } = useLibraryMutation(props.node);
  const ref = useRef<CombotagsHandle>(null);
  const [loadingTags, setLoadingTags] = useState(false);

  const currentTags = props.node.tags.map((t) => t.name);

  async function handleQuery(q: string): Promise<TagNameList> {
    const tags =
      (await handle(async () => {
        const { tags } = await tagList({ q });
        return tags.map((t) => t.name);
      })) ?? [];

    const filtered = tags.filter((t) => !currentTags.includes(t));

    return filtered;
  }

  async function handleSuggestTags() {
    await handle(
      async () => {
        setLoadingTags(true);
        const tags = await suggestTags(props.node.slug);

        if (!tags) {
          throw new Error(
            "No tags could be suggested for this page. This may be due to the content being too short.",
          );
        }

        ref.current?.append(tags);
      },
      {
        async cleanup() {
          setLoadingTags(false);
        },
      },
    );
  }

  return {
    ref,
    currentTags,
    isSuggestEnabled,
    loadingTags,
    handleQuery,
    handleSuggestTags,
  };
}

export function LibraryPageTagsList<T extends FieldValues>(props: Props<T>) {
  const {
    ref,
    currentTags,
    isSuggestEnabled,
    loadingTags,
    handleQuery,
    handleSuggestTags,
  } = useLibraryPageTagsList(props);

  if (props.editing) {
    return (
      <Controller
        name={props.name}
        control={props.control}
        render={({ field }) => {
          async function handleChange(values: string[]) {
            field.onChange(values);
          }

          return (
            <HStack w="full" gap="1" alignItems="start">
              <Combotags
                ref={ref}
                initialValue={currentTags}
                onQuery={handleQuery}
                onChange={handleChange}
              />
              {isSuggestEnabled && (
                <IntelligenceAction
                  title="Suggest tags for this page"
                  onClick={handleSuggestTags}
                  variant="subtle"
                  h="full"
                  loading={loadingTags}
                />
              )}
            </HStack>
          );
        }}
      />
    );
  }

  if (props.node.tags.length === 0) {
    return null;
  }

  return <TagBadgeList tags={props.node.tags} />;
}
