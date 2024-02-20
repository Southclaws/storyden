import { Portal } from "@ark-ui/react";
import { format } from "date-fns/format";

import { DeleteConfirmation } from "src/components/site/Action/Delete";
import { MoreAction } from "src/components/site/Action/More";
import {
  Menu,
  MenuContent,
  MenuItem,
  MenuItemGroup,
  MenuItemGroupLabel,
  MenuPositioner,
  MenuSeparator,
  MenuTrigger,
} from "src/theme/components/Menu";

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
    <Menu size="sm" lazyMount onSelect={handleSelect}>
      <MenuTrigger asChild>
        <MoreAction size="xs" />
      </MenuTrigger>
      <Portal>
        <MenuPositioner>
          <MenuContent minW="36">
            <MenuItemGroup id="user">
              <MenuItemGroupLabel
                htmlFor="user"
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
              </MenuItemGroupLabel>

              <MenuSeparator />

              {reviewFlow && (
                <>
                  {reviewFlow.draftToReview && (
                    <MenuItem id="review">Submit for review</MenuItem>
                  )}
                  {reviewFlow.reviewToPublsh && (
                    <MenuItem id="publsh">Publish</MenuItem>
                  )}
                  {reviewFlow.publishToReview && (
                    <MenuItem id="review">Unpublish</MenuItem>
                  )}
                  {reviewFlow.reviewToDraft && (
                    <MenuItem id="draft">Revert to draft</MenuItem>
                  )}
                  {reviewFlow.draftToPublish && (
                    <MenuItem id="publish">Force publish</MenuItem>
                  )}
                </>
              )}

              {deleteEnabled && (
                <>
                  <MenuItem id="delete">Delete</MenuItem>
                  <DeleteConfirmation {...deleteProps} />
                </>
              )}
            </MenuItemGroup>
          </MenuContent>
        </MenuPositioner>
      </Portal>
    </Menu>
  );
}
