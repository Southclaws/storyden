import { Suspense, lazy } from "react";

import { TagNameList, TagReferenceList } from "@/api/openapi-schema";
import { TagBadgeList } from "@/components/tag/TagBadgeList";

const ThreadTagListEditing = lazy(() =>
  import("./ThreadTagListEditing").then((module) => ({
    default: module.ThreadTagListEditing,
  })),
);

export type Props = {
  editing: boolean;
  disabled?: boolean;
  initialTags?: TagReferenceList;
  onChange: (tags: TagNameList) => Promise<void>;
};

export function ThreadTagList(props: Props) {
  const tags = <TagBadgeList tags={props.initialTags ?? []} />;

  if (props.editing) {
    return (
      <Suspense fallback={tags}>
        <ThreadTagListEditing {...props} />
      </Suspense>
    );
  }

  if (props.initialTags?.length === 0) {
    return null;
  }

  return tags;
}
