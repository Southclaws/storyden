import { motion } from "framer-motion";
import { CircleDashed } from "lucide-react";
import { ComponentProps } from "react";

import { styled } from "@/styled-system/jsx";

const CircleDashedIcon = styled(CircleDashed);

type Props = {
  size?: number;
} & ComponentProps<typeof CircleDashedIcon>;

export function RunningAnimatedIcon({ size = 16, ...rest }: Props) {
  return (
    <motion.div
      animate={{ rotate: 360 }}
      transition={{
        duration: 3,
        repeat: Infinity,
        ease: "linear",
      }}
      style={{ display: "flex", alignItems: "center" }}
    >
      <CircleDashedIcon size={size} {...rest} />
    </motion.div>
  );
}
