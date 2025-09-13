import { usePopover } from "@ark-ui/react";
import { AnimatePresence, motion } from "framer-motion";
import { useState } from "react";
import { toast } from "sonner";
import { match } from "ts-pattern";

import { Spinner } from "@/components/ui/Spinner";
import { Button, ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CheckCircleIcon } from "@/components/ui/icons/CheckCircle";
import { LinkIcon } from "@/components/ui/icons/Link";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { Input } from "@/components/ui/input";
import { Item } from "@/components/ui/menu";
import * as Popover from "@/components/ui/popover";
import { useCapability } from "@/lib/settings/capabilities";
import { HStack, styled } from "@/styled-system/jsx";
import { UtilityValues } from "@/styled-system/types/prop-type";
import { deriveError } from "@/utils/error";

import {
  ImportState,
  ImportStep,
  importFromURLGenerator,
  importStateLabel,
} from "./import";

const POPOVER_CLOSE_DELAY = 3000;

type Props = ButtonProps & {
  parentSlug?: string;
  hideLabel?: boolean;
  onComplete?: () => void;
};

export const CreatePageFromURLID = "create-page-from-url";
export const CreatePageFromURLLabel = "Create from URL";
export const CreatePageFromURLIcon = <LinkIcon />;

export function CreatePageFromURLAction({
  parentSlug,
  hideLabel,
  onComplete,
  ...props
}: Props) {
  const genaiAvailable = useCapability("gen_ai");
  const [url, setUrl] = useState({
    valid: false,
    value: "",
    url: null as URL | null,
  });
  const [importState, setImportState] = useState<ImportState | null>(null);
  const [isImporting, setIsImporting] = useState(false);

  function resetState() {
    setUrl({
      valid: false,
      value: "",
      url: null,
    });
    setImportState(null);
    setIsImporting(false);
  }

  const popover = usePopover({
    onOpenChange: (open) => {
      if (!open) {
        resetState();
      }
    },
  });

  function handleInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    const inputUrl = e.target.value;
    try {
      const parsedUrl = new URL(inputUrl);
      setUrl({
        valid: true,
        value: inputUrl,
        url: parsedUrl,
      });
    } catch {
      setUrl({
        valid: false,
        value: inputUrl,
        url: null,
      });
    }
  }

  async function handleImport() {
    if (!url.valid || !url.url) return;

    setIsImporting(true);
    setImportState(null);

    try {
      const generator = importFromURLGenerator({
        url: url.value,
        parentSlug,
        genaiAvailable,
      });

      for await (const state of generator) {
        if (state.step === "failed") {
          throw new Error(state.error || "Import failed");
        }

        setImportState(state);

        if (state.step === "complete" && !state.error) {
          onComplete?.();
          await new Promise((resolve) =>
            setTimeout(resolve, POPOVER_CLOSE_DELAY),
          );
          popover.setOpen(false);
          resetState();
          break;
        }

        if (state.error) {
          break;
        }
      }
    } catch (error) {
      const derived = deriveError(error);
      toast.error(`Failed to import from URL: ${derived}`);
      setImportState({
        step: "failed",
        error: derived,
      });

      await new Promise((resolve) => setTimeout(resolve, POPOVER_CLOSE_DELAY));
      setImportState(null);
    } finally {
      setIsImporting(false);
    }
  }

  return (
    <Popover.RootProvider value={popover}>
      <Popover.Trigger asChild>
        <IconButton
          type="button"
          size="xs"
          variant="subtle"
          px={hideLabel ? "0" : "1"}
          {...props}
        >
          {CreatePageFromURLIcon}
          {!hideLabel && (
            <>
              <span>{CreatePageFromURLLabel}</span>
            </>
          )}
        </IconButton>
      </Popover.Trigger>
      <Popover.Positioner>
        <Popover.Content px="1" py="1">
          {!importState ? (
            <HStack gap="2" transition="all">
              <Input
                w="64"
                size="xs"
                placeholder="Enter URL to import..."
                value={url.value}
                onChange={handleInputChange}
              />
              <Button
                size="xs"
                onClick={handleImport}
                disabled={!url.valid || isImporting}
                loading={isImporting}
              >
                Import
              </Button>
            </HStack>
          ) : (
            <HStack gap="2" justify="space-between">
              {match(importState.step)
                .with("complete", () => (
                  <CheckCircleIcon color="fg.success" fill="bg.success" />
                ))
                .with("failed", () => (
                  <WarningIcon color="fg.warning" fill="bg.warning" />
                ))
                .otherwise(() => (
                  <Spinner />
                ))}
              <AnimatePresence mode="wait">
                <motion.div
                  key={importState.step}
                  initial={{ opacity: 0, y: 2 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -2 }}
                  transition={{ duration: 0.15 }}
                >
                  <styled.span
                    color={getImportStateColor(importState.step)}
                    px="2"
                  >
                    {importStateLabel[importState.step]}
                  </styled.span>
                </motion.div>
              </AnimatePresence>
            </HStack>
          )}
        </Popover.Content>
      </Popover.Positioner>
    </Popover.RootProvider>
  );
}

export function CreatePageFromURLMenuItem({ hideLabel }: Props) {
  return (
    <Item value={CreatePageFromURLID}>
      {CreatePageFromURLIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{CreatePageFromURLLabel}</span>
        </>
      )}
    </Item>
  );
}

function getImportStateColor(state: ImportStep): UtilityValues["color"] {
  switch (state) {
    case "failed":
      return "fg.warning";
    default:
      return "fg.muted";
  }
}
