import { SVGProps } from "react";

export function SendIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={24}
      height={25}
      fill="none"
      {...props}
    >
      <g
        stroke="#303030"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={1.5}
        clipPath="url(#a)"
      >
        <path d="M4.75 19.75 12 5.25l7.25 14.5-7.25-3.5-7.25 3.5ZM12 16v-2.75" />
      </g>
      <defs>
        <clipPath id="a">
          <rect width={24} height={24} y={0.5} fill="#fff" rx={12} />
        </clipPath>
      </defs>
    </svg>
  );
}
