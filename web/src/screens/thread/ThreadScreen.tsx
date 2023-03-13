import { useRouter } from "next/router";
import { useThreadGet } from "src/api/openapi/threads";
import { Unready } from "src/components/Unready";
import { ThreadView } from "./components/ThreadView/ThreadView";

export function ThreadScreen() {
  const router = useRouter();
  const slug = router.query["slug"] as string;

  const thread = useThreadGet(slug);

  if (!thread.data) return <Unready {...thread.error} />;

  return <ThreadView {...thread.data} />;
}
