import { Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

import { DeleteConfirmation } from "src/components/site/Action/Delete";
import { MoreAction } from "src/components/site/Action/More";

import * as Menu from "@/components/ui/menu";
import { styled } from "@/styled-system/jsx";

import { Props, useDatagraphNodeMenu } from "./useDatagraphNodeMenu";

export function DatagraphNodeMenu(props: Props) {
  const { reviewFlow, deleteEnabled, deleteProps, handleSelect } =
    useDatagraphNodeMenu(props);

  const statusText =
    props.node.visibility === "draft"
      ? "(draft)"
      : props.node.visibility === "review"
        ? "(in review)"
        : "";

  return (
    <Menu.Root size="sm" lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <MoreAction size="xs" />
      </Menu.Trigger>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="36">
            <Menu.ItemGroup id="user">
              <Menu.ItemGroupLabel
                display="flex"
                flexDir="column"
                userSelect="none"
              >
                <styled.span>
                  {`Created by ${props.node.owner.name}`} {statusText}
                </styled.span>

                <styled.time fontWeight="normal">
                  {format(new Date(props.node.createdAt), "yyyy-mm-dd")}
                </styled.time>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              {reviewFlow && (
                <>
                  {reviewFlow.draftToReview && (
                    <Menu.Item value="review">Submit for review</Menu.Item>
                  )}
                  {reviewFlow.reviewToPublsh && (
                    <Menu.Item value="publsh">Publish</Menu.Item>
                  )}
                  {reviewFlow.publishToReview && (
                    <Menu.Item value="review">Unpublish</Menu.Item>
                  )}
                  {reviewFlow.reviewToDraft && (
                    <Menu.Item value="draft">Revert to draft</Menu.Item>
                  )}
                  {reviewFlow.draftToPublish && (
                    <Menu.Item value="publish">Force publish</Menu.Item>
                  )}
                </>
              )}

              {deleteEnabled && (
                <>
                  <Menu.Item value="delete">Delete</Menu.Item>
                  <DeleteConfirmation {...deleteProps} />
                </>
              )}
            </Menu.ItemGroup>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
