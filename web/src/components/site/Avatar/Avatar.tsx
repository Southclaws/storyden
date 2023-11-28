import Image from "next/image";

import { css } from "@/styled-system/css";
import { Box, BoxProps } from "@/styled-system/jsx";

import { useAvatar } from "./useAvatar";

type Props = {
  handle: string;
} & BoxProps;

export function Avatar({ handle, ...props }: Props) {
  const { src } = useAvatar(handle);
  return (
    <Box width="8" {...props}>
      <Image
        className={css({
          borderRadius: "full",
        })}
        src={src}
        width={32}
        height={32}
        alt={`${handle}'s avatar`}
      />
    </Box>
  );
}
