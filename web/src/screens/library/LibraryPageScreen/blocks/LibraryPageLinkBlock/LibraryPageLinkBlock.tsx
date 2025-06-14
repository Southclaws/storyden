import { match } from "ts-pattern";

import { LinkCard } from "@/components/library/links/LinkCard";
import { InfoTip } from "@/components/site/InfoTip";
import { Unready } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { LinkButton } from "@/components/ui/link-button";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useEditState } from "../../useEditState";

import { useLibraryPageLinkBlock } from "./useLibraryPageLinkBlock";

export function LibraryPageLinkBlock() {
  const { editing } = useEditState();
  const { currentNode } = useLibraryPageContext();

  if (editing) {
    return <LibraryPageLinkBlockEditing />;
  }

  if (!currentNode.link?.url) {
    return null;
  }

  return (
    <LinkButton href={currentNode.link.url} size="xs" variant="subtle">
      {currentNode.link?.domain}
    </LinkButton>
  );
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
        .with(null, () => null)
        .with(undefined, () => <Unready />)
        .otherwise((link) => (
          <LinkCard link={link} />
        ))}
    </LStack>
  );
}
