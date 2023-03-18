import { ReplyBox } from "src/components/ReplyBox";
import { useComposeScreen } from "./useComposeScreen";

export function ComposeScreen() {
  const { onCreate } = useComposeScreen();
  return <ReplyBox onSave={onCreate} />;
}
