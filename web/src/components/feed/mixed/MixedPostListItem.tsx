import { ThreadReference } from "src/api/openapi/schemas";

import { LinkPost } from "../link/LinkPost";
import { TextPost } from "../text/TextPost";

import { getPostType } from "./utils";

type Props = {
  thread: ThreadReference;
  onDelete: () => void;
};

export function MixedPostListItem(props: Props) {
  switch (getPostType(props.thread)) {
    case "link":
      return <LinkPost {...props} />;
    default:
      return <TextPost {...props} />;
  }
}
