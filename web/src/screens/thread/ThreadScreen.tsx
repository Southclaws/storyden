import { Unready } from "src/components/Unready";
import { ThreadView } from "./components/ThreadView/ThreadView";
import { ThreadScreenContext } from "./context";
import { useThreadScreen } from "./useThreadScreen";
import { Skeleton, SkeletonText } from "@chakra-ui/react";

export function ThreadScreen() {
  const { state, data, error } = useThreadScreen();

  if (!data)
    return (
      <Unready {...error}>
        <Skeleton height={8} />
        <Skeleton height={4} />
        <Skeleton height={5} />
        <SkeletonText noOfLines={3} />
        <SkeletonText noOfLines={5} />
        <SkeletonText noOfLines={9} />
      </Unready>
    );

  return (
    <ThreadScreenContext.Provider value={state}>
      <ThreadView {...data} />
    </ThreadScreenContext.Provider>
  );
}
