import { EllipsisHorizontalIcon } from "@heroicons/react/24/solid";

import { Button } from "src/theme/components/Button";
import { useDisclosure } from "src/utils/useDisclosure";

import { Box } from "@/styled-system/jsx";

import { MemberOptionsMenu } from "./MemberOptionsMenu";
import { MemberOptionsModal } from "./MemberOptionsModal";
import { Props } from "./useMemberOptionsScreen";

export function MemberOptionsTrigger(props: Props) {
  const { isOpen, onClose, onToggle } = useDisclosure();

  return (
    <Box id="member-options-trigger">
      <Box display={{ base: "none", md: "block" }}>
        <MemberOptionsMenu {...props}>
          <Button size="xs" kind="ghost">
            <EllipsisHorizontalIcon />
          </Button>
        </MemberOptionsMenu>
      </Box>

      <Box display={{ base: "block", md: "none" }}>
        <Button size="xs" kind="ghost" onClick={onToggle}>
          <EllipsisHorizontalIcon />
        </Button>

        <MemberOptionsModal isOpen={isOpen} onClose={onClose} {...props} />
      </Box>
    </Box>
  );
}
