import { match } from "ts-pattern";

import { LinkCard } from "@/components/library/links/LinkCard";
import { InfoTip } from "@/components/site/InfoTip";
import { Unready } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Center, HStack, LStack, WStack } from "@/styled-system/jsx";

import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

import { useLibraryPageLinkBlock } from "./useLibraryPageLinkBlock";

export function LibraryPageLinkBlock() {
  const { editing } = useEditState();

  const link = useWatch((s) => s.draft.link);

  if (editing) {
    return <LibraryPageLinkBlockEditing />;
  }

  if (!link?.url) {
    return null;
  }

  return <LinkCard link={link} />;
}

function LibraryPageLinkBlockEditing() {
  const { data, handlers } = useLibraryPageLinkBlock();

  return (
    <LStack gap="0">
      <WStack>
        <Input
          w="full"
          size="sm"
          variant="ghost"
          color="fg.muted"
          placeholder="External URL..."
          onChange={handlers.handleInputValueChange}
          value={data.inputValue}
          defaultValue={data.defaultLinkURL}
        />

        <HStack>
          <InfoTip title="Generating a page from a URL">
            Importing a URL will fetch the content and store it in this page.
          </InfoTip>
          <Button
            type="button"
            size="xs"
            variant="subtle"
            disabled={!data.resolvedLink}
            loading={data.isImporting}
            onClick={handlers.handleImport}
          >
            Import
          </Button>
        </HStack>
      </WStack>

      {match(data.resolvedLink)
        .with(undefined, () => null)
        .with(null, () => (
          <Center w="full" h="24">
            <Unready />
          </Center>
        ))
        .otherwise((link) => (
          <LinkCard link={link} />
        ))}
    </LStack>
  );
}
