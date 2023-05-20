import { Image, ImageProps } from "@chakra-ui/react";
import { useAvatar } from "./useAvatar";

type Props = {
  handle: string;
} & ImageProps;

export function Avatar({ handle, ...props }: Props) {
  const { src } = useAvatar(handle);
  return (
    <Image
      borderRadius="full"
      boxSize={6}
      src={src}
      alt={`${handle}'s avatar`}
      {...props}
    />
  );
}
