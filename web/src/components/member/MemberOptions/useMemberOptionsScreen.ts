import { PublicProfile } from "src/api/openapi/schemas";

export type Props = PublicProfile & {
  onChange?: () => void;
};

export function useMemberOptionsScreen(props: Props) {
  return props;
}
