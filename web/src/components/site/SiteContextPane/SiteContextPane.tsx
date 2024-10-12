"use client";

import Image from "next/image";

import { ContentComposer } from "@/components/content/ContentComposer/ContentComposer";
import { ContentFormField } from "@/components/content/ContentComposer/ContentField";
import { FormControl } from "@/components/ui/FormControl";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";
import { css } from "@/styled-system/css";
import { Divider, HStack, LStack } from "@/styled-system/jsx";
import { LStack as lstack } from "@/styled-system/patterns";

import { EditAction } from "../Action/Edit";
import { SaveAction } from "../Action/Save";
import { AdminAnchor } from "../Navigation/Anchors/Admin";
import { Unready } from "../Unready";

import { Form, Props, useSiteContextPane } from "./useSiteContextPane";

export function SiteContextPane(props: Props) {
  const { ready, error, form, data, handlers } = useSiteContextPane(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { info, iconURL, isEditingEnabled, isAdmin, editing } = data;

  const isEditingSettings = editing === "settings";

  return (
    <form
      className={lstack({ gap: "1" })}
      onSubmit={handlers.handleSaveSettings}
    >
      <HStack w="full" justify="space-between" alignItems="start">
        {isEditingSettings ? (
          <FormControl>
            <Input {...form.register("title")} />
            <FormErrorText>
              {form.formState.errors.title?.message}
            </FormErrorText>
          </FormControl>
        ) : (
          <Heading textWrap="wrap">{info.title}</Heading>
        )}

        <Image
          className={css({
            borderRadius: "md",
            cursor: isEditingSettings ? "help" : "default",
          })}
          alt="Icon"
          src={iconURL}
          width={32}
          height={32}
          title={
            isEditingSettings
              ? "You can change your community's icon in the admin settings page."
              : undefined
          }
        />
      </HStack>

      {isEditingSettings ? (
        <FormControl>
          <Input size="xs" {...form.register("description")} />
          <FormErrorText>
            {form.formState.errors.description?.message}
          </FormErrorText>
        </FormControl>
      ) : (
        <p>{info.description}</p>
      )}

      {isEditingSettings ? (
        <FormControl>
          <ContentFormField<Form>
            control={form.control}
            name="content"
            initialValue={info.content}
            placeholder="About your community..."
          />
          <FormErrorText>
            {form.formState.errors.content?.message}
          </FormErrorText>
        </FormControl>
      ) : (
        info.content && <ContentComposer initialValue={info.content} disabled />
      )}

      {isEditingEnabled && (
        <LStack>
          <FormErrorText>{form.formState.errors.root?.message}</FormErrorText>

          <Divider />

          <HStack w="full" justify="space-between">
            {isEditingSettings ? (
              <SaveAction type="submit">Save</SaveAction>
            ) : (
              <EditAction onClick={handlers.handleEnableEditing}>
                Edit
              </EditAction>
            )}
            {isAdmin && <AdminAnchor />}
          </HStack>
        </LStack>
      )}
    </form>
  );
}
