import { Box, BoxProps, Image } from "src/theme/components";

import { useAvatar } from "./useAvatar";

type Props = {
  handle: string;
} & BoxProps;

export function Avatar({ handle, ...props }: Props) {
  const { src } = useAvatar(handle);
  return (
    <Box width={6} {...props}>
      <Image borderRadius="full" src={src} alt={`${handle}'s avatar`} />
    </Box>
  );
}
