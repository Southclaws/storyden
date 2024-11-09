import { CubeTransparentIcon } from "@heroicons/react/24/outline";
import { Origami } from "lucide-react";

import { styled } from "@/styled-system/jsx";

export const EmptyIcon = styled(CubeTransparentIcon, {
  base: {
    width: "4",
    color: "fg.subtle",
  },
});

export const EmptyThreadsIcon = styled(Origami, {
  base: {
    width: "4",
    color: "fg.subtle",
  },
});
