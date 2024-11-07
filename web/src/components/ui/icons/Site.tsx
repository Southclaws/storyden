import Image from "next/image";

import { Styles, css } from "@/styled-system/css";
import { getIconURL } from "@/utils/icon";

export function SiteIcon(props: Styles) {
  const src = getIconURL("512x512");

  const imageStyles = css(props);

  return (
    <Image className={imageStyles} src={src} alt="" width={512} height={512} />
  );
}
