import { ThreadReference } from "@/api/openapi-schema";
import { Permission } from "@/api/openapi-schema";
import { PermissionBadge } from "@/components/role/PermissionBadge";
import { Timestamp } from "@/components/site/Timestamp";
import * as Alert from "@/components/ui/alert";
import { WarningIcon } from "@/components/ui/icons/Warning";

type Props = {
  thread: ThreadReference;
};

export function ThreadDeletedAlert({ thread }: Props) {
  if (!thread.deletedAt) {
    return null;
  }

  return (
    <Alert.Root
      mb="2"
      colorPalette="red"
      backgroundColor="colorPalette.5"
      borderColor="colorPalette.3"
    >
      <Alert.Icon asChild>
        <WarningIcon />
      </Alert.Icon>
      <Alert.Content>
        <Alert.Title>Thread deleted</Alert.Title>
        <Alert.Description>
          This thread was deleted <Timestamp created={thread.deletedAt} /> and
          is not accessible to any member except the author and members with{" "}
          <PermissionBadge permission={Permission.MANAGE_POSTS} />
        </Alert.Description>
      </Alert.Content>
    </Alert.Root>
  );
}
