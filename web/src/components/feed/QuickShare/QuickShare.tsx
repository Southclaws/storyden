"use client";

import { Suspense, lazy, useState } from "react";

import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";

import type { Props } from "./useQuickShare";

const QuickShareForm = lazy(() =>
  import("./QuickShareForm").then((module) => ({
    default: module.QuickShareForm,
  })),
);

export function QuickShare(props: Props) {
  const [isActive, setIsActive] = useState(false);
  const [focusOnReady, setFocusOnReady] = useState(false);

  if (!props.initialSession) {
    return null;
  }

  function handleActivate() {
    setFocusOnReady(true);
    setIsActive(true);
  }

  const prompt = (
    <Button
      className="quick-share__trigger"
      type="button"
      variant="ghost"
      aria-label="Quickly create a new thread"
      aria-busy={isActive || undefined}
      w="full"
      minH="16"
      px="0"
      alignItems="start"
      justifyContent="start"
      fontWeight="normal"
      fontSize="md"
      color="text.muted"
      textAlign="start"
      whiteSpace="normal"
      onFocus={handleActivate}
      onClick={handleActivate}
    >
      Share a thought, a link, something cool...
    </Button>
  );

  return (
    <CardBox
      className="quick-share"
      bgColor="background.control"
      borderRadius="md"
      onPointerEnter={(event) => {
        if (event.pointerType === "mouse") {
          setIsActive(true);
        }
      }}
    >
      <Suspense fallback={prompt}>
        {isActive ? (
          <QuickShareForm {...props} focusOnReady={focusOnReady} />
        ) : (
          prompt
        )}
      </Suspense>
    </CardBox>
  );
}
