import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import { Node, TagNameList } from "@/api/openapi-schema";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Combotags } from "@/components/ui/combotags";
import { useLibraryMutation } from "@/lib/library/library";

export type Props = {
  editing: boolean;
  node: Node;
};

export function LibraryPageTagsList(props: Props) {
  const { updateNode, revalidate } = useLibraryMutation(props.node);

  const currentTags = props.node.tags.map((t) => t.name);

  async function handleChange(values: string[]) {
    await handle(
      async () => {
        await updateNode(props.node.slug, { tags: values });
      },
      {
        cleanup: async () => await revalidate(),
      },
    );
  }

  async function handleQuery(q: string): Promise<TagNameList> {
    const tags =
      (await handle(async () => {
        const { tags } = await tagList({ q });
        return tags.map((t) => t.name);
      })) ?? [];

    const filtered = tags.filter((t) => !currentTags.includes(t));

    return filtered;
  }

  if (props.editing) {
    return (
      <>
        <Combotags
          initialValue={currentTags}
          onQuery={handleQuery}
          onChange={handleChange}
        />
      </>
    );
  }

  if (props.node.tags.length === 0) {
    return null;
  }

  return <TagBadgeList tags={props.node.tags} />;
}
