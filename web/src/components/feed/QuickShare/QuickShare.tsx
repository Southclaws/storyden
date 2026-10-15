"use client";

import { useState } from "react";

import { CardBox } from "@/components/ui/card-box";

import { QuickShareForm } from "./QuickShareForm";
import type { Props } from "./useQuickShare";

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
      <QuickShareForm
        {...props}
        isActive={isActive}
        focusOnReady={focusOnReady}
        onActivate={handleActivate}
      />
    </CardBox>
  );
}
