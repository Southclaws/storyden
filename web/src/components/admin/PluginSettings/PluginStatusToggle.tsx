import { AnimatePresence, motion } from "framer-motion";
import { PauseCircleIcon } from "lucide-react";
import { useState } from "react";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  getPluginGetKey,
  usePluginSetActiveState,
} from "@/api/openapi-client/plugins";
import { Plugin, PluginActiveState } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { RunningAnimatedIcon } from "@/components/ui/icons/RunningAnimatedIcon";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { HStack } from "@/styled-system/jsx";

import { isPluginStatusActive } from "./utils";

type Props = { plugin: Plugin };

export function PluginStatusToggle({ plugin }: Props) {
  const { mutate } = useSWRConfig();
  const [transitioningTo, setTransitioningTo] =
    useState<PluginActiveState | null>(null);
  const { trigger: setActiveState } = usePluginSetActiveState(plugin.id);

  const activeState = plugin.status.active_state;
  const isActive = isPluginStatusActive(plugin.status);
  const isTransitioning = transitioningTo !== null;

  const handleToggleActive = async () => {
    if (isTransitioning) return;

    const targetState = isActive
      ? PluginActiveState.inactive
      : PluginActiveState.active;

    setTransitioningTo(targetState);

    const minimumDuration = 1000;
    const startTime = Date.now();

    await handle(
      async () => {
        await setActiveState({ active: targetState });

        const elapsedTime = Date.now() - startTime;
        const remainingTime = Math.max(0, minimumDuration - elapsedTime);

        if (remainingTime > 0) {
          await new Promise((resolve) => setTimeout(resolve, remainingTime));
        }

        setTransitioningTo(null);

        await mutate(getPluginGetKey(plugin.id));
      },
      {
        onError: async () => {
          setTransitioningTo(null);
        },
        cleanup: async () => {
          setTransitioningTo(null);
        },
      },
    );
  };

  const displayState = transitioningTo ?? activeState;
  const statusLabel = isTransitioning
    ? getTransitionLabel(activeState, transitioningTo!)
    : getStatusLabel(displayState);

  const icon = (() => {
    switch (displayState) {
      case PluginActiveState.active:
        return <RunningAnimatedIcon size={16} color="accent.9" />;
      case PluginActiveState.inactive:
        return <PauseCircleIcon size={16} />;
      case PluginActiveState.error:
        return <WarningIcon size={16} />;
      default:
        return null;
    }
  })();

  const actionLabel = (() => {
    switch (activeState) {
      case PluginActiveState.active:
        return "Disable";
      case PluginActiveState.inactive:
        return "Enable";
      case PluginActiveState.error:
        return "Retry";
      default:
        return "Unknown";
    }
  })();

  return (
    <HStack
      style={
        {
          "--btn-min-width": "8ch",
        } as any
      }
      gap="0"
    >
      <Badge
        borderRightRadius="none"
        borderLeftRadius="full"
        px="1"
        overflow="hidden"
        bgColor={isActive ? "accent.4" : "bg.subtle"}
        borderColor={isActive ? "accent.3" : "border.subtle"}
        color={isActive ? "accent.9" : "fg.subtle"}
        borderRightWidth="none"
      >
        <AnimatePresence mode="wait">
          <motion.div
            key={displayState}
            initial={{ scale: 0.8, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.8, opacity: 0 }}
            transition={{ type: "spring", stiffness: 500, damping: 30 }}
            style={{ display: "flex", alignItems: "center" }}
          >
            {icon}
          </motion.div>
        </AnimatePresence>

        <AnimatePresence mode="wait">
          <motion.span
            key={statusLabel}
            initial={{ opacity: 0, y: 5 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -5 }}
            transition={{ duration: 0.2 }}
          >
            {statusLabel}
          </motion.span>
        </AnimatePresence>
      </Badge>

      <Button
        borderLeftRadius="none"
        borderRightRadius="full"
        size="xs"
        minWidth="var(--btn-min-width)"
        onClick={handleToggleActive}
        disabled={isTransitioning}
        variant={isActive ? "subtle" : "solid"}
        bgColor={isActive ? "bg.subtle" : "accent.4"}
        borderWidth="thin"
        borderLeftWidth="none"
        borderColor={isActive ? "border.subtle" : "accent.5"}
        color={isActive ? "fg.subtle" : "accent.9"}
        _hover={{
          bgColor: isActive ? "bg.muted" : "accent.5",
        }}
      >
        {actionLabel}
      </Button>
    </HStack>
  );
}

function getStatusLabel(activeState: PluginActiveState) {
  switch (activeState) {
    case PluginActiveState.active:
      return "Running";
    case PluginActiveState.inactive:
      return "Disabled";
    case PluginActiveState.error:
      return "Error";
    default:
      return "Unknown";
  }
}

function getTransitionLabel(
  fromState: PluginActiveState,
  toState: PluginActiveState,
) {
  if (
    fromState === PluginActiveState.active &&
    toState === PluginActiveState.inactive
  ) {
    return "Disabling";
  }

  if (
    fromState === PluginActiveState.inactive &&
    toState === PluginActiveState.active
  ) {
    return "Enabling";
  }

  return `Pending`;
}
