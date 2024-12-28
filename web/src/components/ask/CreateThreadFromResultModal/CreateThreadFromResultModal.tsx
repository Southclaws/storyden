import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Button } from "@/components/ui/button";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { UseDisclosureProps, useDisclosure } from "@/utils/useDisclosure";

import { CreateThreadFromResultScreen } from "./CreateThreadFromResultScreen";
import { Props } from "./useCreateThreadFromResult";

export function CreateThreadFromResultModal({
  contentMarkdown,
  sources,
  onClose,
  onOpen,
  isOpen,
}: UseDisclosureProps & Props) {
  return (
    <ModalDrawer
      onOpen={onOpen}
      isOpen={isOpen}
      onClose={onClose}
      title={`Post thread from Ask result`}
    >
      <CreateThreadFromResultScreen
        onFinish={onClose}
        contentMarkdown={contentMarkdown}
        sources={sources}
      />
    </ModalDrawer>
  );
}

export function CreateThreadFromResultModalTrigger({
  contentMarkdown,
  sources,
}: Props) {
  const disclosure = useDisclosure();

  return (
    <>
      <Button
        size="xs"
        title="Create a thread from this Ask result"
        onClick={disclosure.onOpen}
      >
        <DiscussionIcon />
        Create Thread
      </Button>

      <CreateThreadFromResultModal
        {...disclosure}
        contentMarkdown={contentMarkdown}
        sources={sources}
      />
    </>
  );
}
