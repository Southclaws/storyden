import { PublicProfile } from "src/api/openapi-schema";

export type Props = PublicProfile & {
  onChange?: () => void;
};

export function useMemberOptionsScreen(props: Props) {
  return props;
}
