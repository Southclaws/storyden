import { TextPostList } from "src/components/feed/text/TextPostList";
import { Unready } from "src/components/site/Unready";

import { useContent } from "./useContent";

export function Content(props: { showEmptyState: boolean }) {
  const { data, error } = useContent();

  if (!data) return <Unready {...error} />;

  return (
    <TextPostList showEmptyState={props.showEmptyState} posts={data.threads} />
  );
}
