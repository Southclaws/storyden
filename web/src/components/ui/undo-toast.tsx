"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";

import { Box, HStack } from "@/styled-system/jsx";

import { Button } from "./button";

type UndoToastProps = {
  message: string;
  duration?: number;
  toastId?: string;
  onUndo: () => void;
  onComplete: () => void;
};

type UndoToastContentProps = Omit<UndoToastProps, "toastId"> & {
  toastId: string | number;
};

function UndoToastContent({
  message,
  duration = 5000,
  onUndo,
  onComplete,
  toastId,
}: UndoToastContentProps) {
  const [progress, setProgress] = useState(100);
  const [isUndone, setIsUndone] = useState(false);
  const hasCompletedRef = useRef(false);

  const complete = useCallback(() => {
    if (hasCompletedRef.current) return;
    hasCompletedRef.current = true;
    onComplete();
  }, [onComplete]);

  useEffect(() => {
    if (isUndone) return;

    const startTime = Date.now();
    let animationFrameId: number;

    const updateProgress = () => {
      const elapsed = Date.now() - startTime;
      const remaining = Math.max(0, 100 - (elapsed / duration) * 100);

      setProgress(remaining);

      if (remaining === 0) {
        toast.dismiss(toastId);
        complete();
      } else {
        animationFrameId = requestAnimationFrame(updateProgress);
      }
    };

    animationFrameId = requestAnimationFrame(updateProgress);

    return () => cancelAnimationFrame(animationFrameId);
  }, [duration, isUndone, toastId, complete]);

  const handleUndo = () => {
    setIsUndone(true);
    onUndo();
    toast.dismiss(toastId);
  };

  const handleClose = () => {
    toast.dismiss(toastId);
    complete();
  };

  return (
    <div
      style={{
        width: "100%",
        display: "flex",
        flexDirection: "column",
        gap: "0.5rem",
      }}
    >
      <HStack justify="space-between" alignItems="center" gap="2">
        <span style={{ fontSize: "0.875rem", flex: 1 }}>{message}</span>
        <HStack gap="2">
          <Button
            size="sm"
            variant="subtle"
            onClick={handleUndo}
            disabled={isUndone}
          >
            {isUndone ? "Cancelled" : "Undo"}
          </Button>
          <Button
            size="sm"
            variant="ghost"
            onClick={handleClose}
          >
            Delete now
          </Button>
        </HStack>
      </HStack>

      <div
        style={{
          width: "100%",
          height: "0.25rem",
          backgroundColor: "var(--colors-bg-muted)",
          borderRadius: "9999px",
          overflow: "hidden",
        }}
      >
        <div
          style={{
            height: "100%",
            width: `${progress}%`,
            borderRadius: "9999px",
            backgroundColor:
              "light-dark(var(--accent-colour-flat-fill-600), var(--accent-colour-dark-fill-600))",
          }}
        />
      </div>
    </div>
  );
}

export function showUndoToast({
  message,
  duration = 5000,
  toastId: customToastId,
  onUndo,
  onComplete,
}: UndoToastProps) {
  const toastId = toast.custom(
    (t) => (
      <Box
        backgroundColor="bg.default"
        borderWidth="thin"
        borderColor="border.default"
        borderRadius="lg"
        padding="3"
        boxShadow="lg"
        minWidth="sm"
        maxWidth="md"
      >
        <UndoToastContent
          message={message}
          duration={duration}
          onUndo={onUndo}
          onComplete={onComplete}
          toastId={t}
        />
      </Box>
    ),
    {
      id: customToastId,
      duration: Infinity,
      position: "bottom-right",
    },
  );

  return toastId;
}
