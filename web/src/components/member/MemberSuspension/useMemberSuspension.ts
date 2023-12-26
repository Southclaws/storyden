import {
  adminAccountBanCreate,
  adminAccountBanRemove,
} from "src/api/openapi/accounts";
import { useProfileGet } from "src/api/openapi/profiles";
import { PublicProfile } from "src/api/openapi/schemas";
import { WithDisclosure } from "src/utils/useDisclosure";

export type Props = PublicProfile & {
  onChange?: () => void;
};

export function useMemberSuspension(props: WithDisclosure<Props>) {
  const { mutate } = useProfileGet(props.handle);

  async function handleSuspension() {
    await adminAccountBanCreate(props.handle);

    mutate();
    props.onChange?.();
    props.onClose?.();
  }

  async function handleReinstate() {
    await adminAccountBanRemove(props.handle);

    mutate();
    props.onChange?.();
    props.onClose?.();
  }

  return {
    handlers: {
      handleSuspension,
      handleReinstate,
    },
  };
}
