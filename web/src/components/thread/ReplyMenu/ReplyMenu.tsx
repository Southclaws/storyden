"use client";

import { Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

import { MoreAction } from "src/components/site/Action/More";

import { DatagraphItemKind } from "@/api/openapi-schema";
import {
  ReportPostMenuItem,
  truncateBody,
} from "@/components/report/ReportPostMenuItem";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { EditIcon } from "@/components/ui/icons/Edit";
import { LinkIcon } from "@/components/ui/icons/Link";
import { ShareIcon } from "@/components/ui/icons/Share";
import * as Menu from "@/components/ui/menu";
import { useI18n } from "@/i18n/provider";
import { HStack, styled } from "@/styled-system/jsx";

import { Props, useReplyMenu } from "./useReplyMenu";

export function ReplyMenu(props: Props) {
  const { t } = useI18n();
  const { isSharingEnabled, isEditingEnabled, isDeletingEnabled, handlers } =
    useReplyMenu(props);

  return (
    <Menu.Root lazyMount>
      <Menu.Trigger asChild>
        <MoreAction size="xs" />
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.ItemGroup id="group">
              <Menu.ItemGroupLabel
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <styled.span>
                  {t("Post by {{name}}", { name: props.reply.author.name })}
                </styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(props.reply.createdAt), "yyyy-MM-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link" onClick={handlers.handleCopyURL}>
                <HStack gap="1">
                  <LinkIcon /> {t("Copy link")}
                </HStack>
              </Menu.Item>

              {isSharingEnabled && (
                <Menu.Item value="share" onClick={handlers.handleShare}>
                  <HStack gap="1">
                    <ShareIcon /> {t("Share")}
                  </HStack>
                </Menu.Item>
              )}

              <ReportPostMenuItem
                menuLabel="Report reply"
                targetKind={DatagraphItemKind.reply}
                targetId={props.reply.id}
                author={props.reply.author}
                headline={t("Reply from {{name}}", {
                  name: props.reply.author.name,
                })}
                body={truncateBody(props.reply.body)}
              />

              {isEditingEnabled && (
                <Menu.Item value="edit" onClick={handlers.handleSetEditing}>
                  <HStack gap="1">
                    <EditIcon /> {t("Edit")}
                  </HStack>
                </Menu.Item>
              )}

              {isDeletingEnabled && (
                <Menu.Item value="delete" onClick={handlers.handleDelete}>
                  <HStack gap="1">
                    <DeleteIcon /> {t("Delete")}
                  </HStack>
                </Menu.Item>
              )}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
