import { useThreadsList } from "src/api/openapi/threads";
import { Unready } from "src/components/Unready";
import { ThreadList } from "./components/ThreadList";

export function HomeScreen() {
  const threads = useThreadsList();

  if (!threads.data) return <Unready {...threads.error} />;

  return <ThreadList threads={threads.data} />;
}
