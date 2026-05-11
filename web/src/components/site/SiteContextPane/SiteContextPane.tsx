"use client";

import Image from "next/image";
import dynamic from "next/dynamic";

import { FormControl } from "@/components/ui/FormControl";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";
import { useI18n } from "@/i18n/provider";
import { css } from "@/styled-system/css";
import { Divider, HStack, LStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { EditAction } from "../Action/Edit";
import { SaveAction } from "../Action/Save";
import { AdminAnchor } from "../Navigation/Anchors/Admin";
import { Unready } from "../Unready";

import { getDisplayContent, getDisplayDescription } from "./defaultContent";
import { Props, useSiteContextPane } from "./useSiteContextPane";

const SiteContextPaneContentField = dynamic(
  () =>
    import("./SiteContextPaneEditor").then(
      (mod) => mod.SiteContextPaneContentField,
    ),
  {
    ssr: false,
  },
);

export function SiteContextPane(props: Props) {
  const { ready, error, form, data, handlers } = useSiteContextPane(props);
  const { locale, t } = useI18n();
  if (!ready) {
    return <Unready error={error} />;
  }

  const { settings, iconURL, isEditingEnabled, isAdmin, editing } = data;

  const isEditingSettings = editing === "settings";
  const displayDescription = getDisplayDescription(
    locale,
    settings.description,
  );
  const displayContent = getDisplayContent(locale, settings.content);

  return (
    <form
      className={lstack({ gap: "1" })}
      onSubmit={handlers.handleSaveSettings}
    >
      <WStack alignItems="start">
        {isEditingSettings ? (
          <FormControl>
            <Input placeholder={t("Site title...")} {...form.register("title")} />
            <FormErrorText>
              {form.formState.errors.title?.message}
            </FormErrorText>
          </FormControl>
        ) : (
          <Heading textWrap="wrap">{settings.title}</Heading>
        )}

        <Image
          className={css({
            borderRadius: "md",
            cursor: isEditingSettings ? "help" : "default",
          })}
          alt={t("Icon")}
          src={iconURL}
          width={32}
          height={32}
          title={
            isEditingSettings
              ? t("You can change your community's icon in the admin settings page.")
              : undefined
          }
        />
      </WStack>

      {isEditingSettings ? (
        <FormControl>
          <Input
            size="xs"
            placeholder={t("Site description...")}
            {...form.register("description")}
          />
          <FormErrorText>
            {form.formState.errors.description?.message}
          </FormErrorText>
        </FormControl>
      ) : (
        <p>{displayDescription}</p>
      )}

      {isEditingSettings ? (
        <FormControl>
          <SiteContextPaneContentField
            control={form.control}
            initialValue={settings.content}
            placeholder={t("About your community...")}
          />
          <FormErrorText>
            {form.formState.errors.content?.message}
          </FormErrorText>
        </FormControl>
      ) : (
        displayContent && (
          <div
            className="typography"
            dangerouslySetInnerHTML={{ __html: displayContent }}
          />
        )
      )}

      {isEditingEnabled && (
        <LStack>
          <FormErrorText>{form.formState.errors.root?.message}</FormErrorText>

          <Divider />

          <WStack>
            {isEditingSettings ? (
              <SaveAction type="submit">{t("Save")}</SaveAction>
            ) : (
              <EditAction onClick={handlers.handleEnableEditing}>
                {t("Edit")}
              </EditAction>
            )}
            {isAdmin && <AdminAnchor />}
          </WStack>
        </LStack>
      )}
    </form>
  );
}
