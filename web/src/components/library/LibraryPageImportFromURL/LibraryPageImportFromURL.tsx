"use client";

import { ChangeEvent, useState } from "react";
import {
  Controller,
  ControllerProps,
  ControllerRenderProps,
} from "react-hook-form";
import { toast } from "sonner";
import { match } from "ts-pattern";

import { handle } from "@/api/client";
import { linkCreate } from "@/api/openapi-client/links";
import { LinkReference, Node } from "@/api/openapi-schema";
import { InfoTip } from "@/components/site/InfoTip";
import { Unready } from "@/components/site/Unready";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Form } from "@/screens/library/LibraryPageScreen/useLibraryPageScreen";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

import { LinkCard } from "../links/LinkCard";

type Props = Omit<ControllerProps<Form>, "render"> & {
  node: Node;
  onImport: (link: LinkReference) => Promise<void>;
};

function useLibraryPageImportFromURL({ node, onImport }: Props) {
  const [link, setLink] = useState<LinkReference | null | undefined>(null);

  async function handleURL(s: string) {
    if (s === "") {
      setLink(null);
      return;
    }

    await handle(async () => {
      try {
        setLink(undefined);

        const u = new URL(s);

        await new Promise((resolve) => setTimeout(resolve, 1000));

        const link = await linkCreate({ url: u.toString() });
        setLink(link);
      } catch (_) {
        // do nothing for invalid URL, already handled by parent form logic.
        setLink(null);
      }
    });
  }

  async function handleImport() {
    if (!link) {
      toast.error("No link available to import.");
      return;
    }

    await onImport(link);
  }

  return {
    data: {
      link,
    },
    handlers: {
      handleURL,
      handleImport,
    },
  };
}

export function LibraryPageImportFromURL(props: Props) {
  const { data, handlers } = useLibraryPageImportFromURL(props);

  const { link } = data;

  return (
    <Controller<Form>
      control={props.control}
      name={props.name}
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
