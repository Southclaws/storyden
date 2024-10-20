"use client";

import { useState } from "react";

export function useConfirmation(fn: () => Promise<void>) {
  const [isConfirming, setConfirming] = useState(false);

  async function handleConfirmAction() {
    if (isConfirming) {
      await fn();
      setConfirming(false);
    } else {
      setConfirming(true);
    }
  }

  async function handleCancelAction() {
    setConfirming(false);
  }

  return {
    isConfirming,
    handleConfirmAction,
    handleCancelAction,
  };
}
