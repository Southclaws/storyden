import { SVGProps } from "react";

export function DragHandleIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={24}
      height={24}
      fill="none"
      {...props}
    >
      <path
        stroke="currentColor"
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12.417 7a.417.417 0 1 1-.834 0 .417.417 0 0 1 .834 0ZM12.417 12a.417.417 0 1 1-.834 0 .417.417 0 0 1 .834 0ZM17.417 7a.417.417 0 1 1-.834 0 .417.417 0 0 1 .834 0ZM17.417 12a.417.417 0 1 1-.834 0 .417.417 0 0 1 .834 0ZM7.417 7a.417.417 0 1 1-.834 0 .417.417 0 0 1 .834 0ZM7.417 12a.417.417 0 1 1-.834 0 .417.417 0 0 1 .834 0ZM12.417 17a.417.417 0 1 1-.834 0 .417.417 0 0 1 .834 0ZM17.417 17a.417.417 0 1 1-.834 0 .417.417 0 0 1 .834 0ZM7.417 17a.417.417 0 1 1-.834 0 .417.417 0 0 1 .834 0Z"
      />
    </svg>
  );
}
