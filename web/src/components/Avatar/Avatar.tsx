import { Image, ImageProps } from "@chakra-ui/react";
import { Account } from "src/api/openapi/schemas";
import { useAvatar } from "./useAvatar";

type Props = {
  account: Account;
} & ImageProps;

export function Avatar({ account, ...props }: Props) {
  const { src, fallback } = useAvatar(account.handle ?? "unknown");
  return (
    <Image
      borderRadius="full"
      boxSize={6}
      src={src}
      fallbackSrc={fallback}
      alt="pic"
      {...props}
    />
  );
}
