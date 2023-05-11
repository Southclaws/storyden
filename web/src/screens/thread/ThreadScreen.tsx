import { Unready } from "src/components/Unready";
import { ThreadView } from "./components/ThreadView/ThreadView";
import { ThreadScreenContext } from "./context";
import { useThreadScreen } from "./useThreadScreen";

export function ThreadScreen() {
  const { state, data, error } = useThreadScreen();

  if (!data) return <Unready {...error} />;

  return (
    <ThreadScreenContext.Provider value={state}>
      <ThreadView {...data} />
    </ThreadScreenContext.Provider>
  );
}
