"use client";

import { Skeleton, SkeletonText } from "@chakra-ui/react";

import { Unready } from "src/components/site/Unready";

import { ThreadView } from "./components/ThreadView/ThreadView";
import { ThreadScreenContext } from "./context";
import { useThreadScreen } from "./useThreadScreen";

// TODO: Make thread view SSR and remove use client from this screen root.

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
