import { usePopover } from "@ark-ui/react";
import { useState } from "react";
import { toast } from "sonner";
import { match } from "ts-pattern";

import { Button, ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CheckCircleIcon } from "@/components/ui/icons/CheckCircle";
import { LinkIcon } from "@/components/ui/icons/Link";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { Input } from "@/components/ui/input";
import { Item } from "@/components/ui/menu";
import * as Popover from "@/components/ui/popover";
import { Spinner } from "@/components/ui/spinner";
import { useCapability } from "@/lib/settings/capabilities";
import { HStack, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

import {
  ImportState,
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
          variant="subtle"
          px={hideLabel ? "0" : "1"}
          aria-label={hideLabel ? CreatePageFromURLLabel : undefined}
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
                placeholder="Enter URL to import..."
                value={url.value}
                onChange={handleInputChange}
              />
              <Button
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
                  <CheckCircleIcon w="5" h="5" color="status.success.content" />
                ))
                .with("failed", () => (
                  <WarningIcon w="5" h="5" color="status.danger.content" />
                ))
                .otherwise(() => (
                  <Spinner />
                ))}
              <styled.div
                key={importState.step}
                animation="[slideFadeIn 150ms ease-out]"
                css={{
                  "@media (prefers-reduced-motion: reduce)": {
                    animation: "[none]",
                  },
                }}
              >
                <styled.span color="text.subtle" px="2">
                  {importStateLabel[importState.step]}
                </styled.span>
              </styled.div>
            </HStack>
          )}
        </Popover.Content>
      </Popover.Positioner>
    </Popover.RootProvider>
  );
}

export function CreatePageFromURLMenuItem({ hideLabel }: Props) {
  return (
    <Item
      value={CreatePageFromURLID}
      aria-label={hideLabel ? CreatePageFromURLLabel : undefined}
    >
      {CreatePageFromURLIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{CreatePageFromURLLabel}</span>
        </>
      )}
    </Item>
  );
}
