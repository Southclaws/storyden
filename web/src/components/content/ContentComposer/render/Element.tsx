import { PropsWithChildren } from "react";
import { RenderElementProps } from "slate-react";

import { RichLink } from "../components/RichLink/RichLink";
import { getURL } from "../utils";

import { styled } from "@/styled-system/jsx";

export function Element({
  children,
  ...props
}: PropsWithChildren<RenderElementProps>) {
  switch (props.element.type) {
    case "paragraph": {
      const url = getURL(props.element);
      if (url) {
        return (
          <RichLink {...props} href={url}>
            {children}
          </RichLink>
        );
      }

      return <styled.p {...props.attributes}>{children}</styled.p>;
    }

    case "link":
      return (
        <a href={props.element.link} {...props.attributes}>
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
      return <styled.ol {...props.attributes}>{children}</styled.ol>;

    case "ul_list":
      return <styled.ul {...props.attributes}>{children}</styled.ul>;

    case "list_item":
      return <styled.li {...props.attributes}>{children}</styled.li>;

    case "image":
      return (
        <>
          <styled.img src={props.element.link} alt="" {...props.attributes} />
          {children}
        </>
      );

    case "block_quote":
      return (
        <styled.blockquote {...props.attributes}>{children}</styled.blockquote>
      );

    case "code_block":
      return (
        <styled.pre overflowX="scroll" width="full" maxW="40">
          {children}
        </styled.pre>
      );

    case "thematic_break":
      return <hr />;

    default:
      console.error("Unknown markdown element rendered", props.element);
      return null;
  }
}
