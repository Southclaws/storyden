import { styled } from "@/styled-system/jsx";
import { ComponentProps, JsxHTMLProps } from "@/styled-system/types";

export function DotSeparator(props: JsxHTMLProps<ComponentProps<"span">>) {
  return (
    <styled.span mx="1" fontWeight="bold" color="fg.subtle" {...props}>
      •
    </styled.span>
  );
}
