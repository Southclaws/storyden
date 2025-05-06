import { ChangeEvent } from "react";
import { Controller, ControllerRenderProps } from "react-hook-form";
import { match } from "ts-pattern";

import { LinkCard } from "@/components/library/links/LinkCard";
import { InfoTip } from "@/components/site/InfoTip";
import { Unready } from "@/components/site/Unready";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { LinkButton } from "@/components/ui/link-button";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { Form } from "../../form";
import { useEditState } from "../../useEditState";

import { useLibraryPageLinkBlock } from "./useLibraryPageLinkBlock";

export function LibraryPageLinkBlock() {
  const { editing } = useEditState();
  const { node } = useLibraryPageContext();

  if (editing) {
    return <LibraryPageLinkBlockEditing />;
  }

  if (!node.link?.url) {
    return null;
  }

  return (
    <LinkButton href={node.link.url} size="xs" variant="subtle">
      {node.link?.domain}
    </LinkButton>
  );
}

function LibraryPageLinkBlockEditing() {
  const { form, data, handlers } = useLibraryPageLinkBlock();

  const { link, isImporting } = data;

  return (
    <Controller<Form>
      control={form.control}
      name="link"
      render={(form) => {
        function handleChange(e: ChangeEvent<HTMLInputElement>) {
          handlers.handleURL(e.target.value);
          form.field.onChange(e);
        }

        const value = form.field.value as ControllerRenderProps<
          Form,
          "link"
        >["value"];

        return (
          <LStack gap="0">
            <WStack>
              <Input
                w="full"
                size="sm"
                variant="ghost"
                color="fg.muted"
                placeholder="External URL..."
                onChange={handleChange}
                value={value}
              />

              <HStack>
                <InfoTip title="Generating a page from a URL">
                  Importing a URL will fetch the content and store it in this
                  page.
                </InfoTip>
                <Button
                  type="button"
                  size="xs"
                  variant="subtle"
                  disabled={!link}
                  loading={isImporting}
                  onClick={handlers.handleImport}
                >
                  Import
                </Button>
              </HStack>
            </WStack>
            <FormErrorText>{form.fieldState.error?.message}</FormErrorText>
            {match(link)
              .with(null, () => null)
              .with(undefined, () => <Unready />)
              .otherwise((link) => (
                <LinkCard link={link} />
              ))}
          </LStack>
        );
      }}
    />
  );
}
