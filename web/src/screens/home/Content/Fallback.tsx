import { ThreadList } from "src/api/openapi/schemas";

import { ThreadList as ThreadListComponent } from "../components/ThreadList";

export function Fallback(props: { threads: ThreadList }) {
  return <ThreadListComponent showEmptyState={true} threads={props.threads} />;
}
