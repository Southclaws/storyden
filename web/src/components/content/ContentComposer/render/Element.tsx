import { PropsWithChildren } from "react";
import { RenderElementProps } from "slate-react";

import { styled } from "@/styled-system/jsx";

export function Element({
  attributes,
  children,
  element,
}: PropsWithChildren<RenderElementProps>) {
  switch (element.type) {
    case "paragraph":
      return <styled.p {...attributes}>{children}</styled.p>;

    case "link":
      return (
        <a href={element.link} {...attributes}>
          {children}
        </a>
      );

    case "heading_one":
      return <h1>{children}</h1>;

    case "heading_two":
      return <h2>{children}</h2>;

    case "heading_three":
      return <h3>{children}</h3>;

    case "heading_four":
      return <h4>{children}</h4>;

    case "heading_five":
      return <h5>{children}</h5>;

    case "heading_six":
      return <h6>{children}</h6>;

    case "ol_list":
      return <styled.ol {...attributes}>{children}</styled.ol>;

    case "ul_list":
      return <styled.ul {...attributes}>{children}</styled.ul>;

    case "list_item":
      return <styled.li {...attributes}>{children}</styled.li>;

    case "image":
      return (
        <>
          <styled.img src={element.link} alt="" {...attributes} />
          {children}
        </>
      );

    case "block_quote":
      return <styled.blockquote {...attributes}>{children}</styled.blockquote>;

    case "code_block":
      return (
        <styled.pre overflowX="scroll" width="full" maxW="40">
          {children}
        </styled.pre>
      );

    case "thematic_break":
      return <hr />;

    default:
      console.error("Unknown markdown element rendered", element);
      return null;
  }
}
