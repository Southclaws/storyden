import { SVGProps } from "react";

export function LoginIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={24}
      height={24}
      fill="none"
      {...props}
    >
      <path
        stroke="#141414"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={1.5}
        d="m9.75 8.75 3.5 3.25-3.5 3.25"
      />
      <path
        stroke="#141414"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={1.5}
        d="M9.75 4.75h7.5a2 2 0 0 1 2 2v10.5a2 2 0 0 1-2 2h-7.5M13 12H4.75"
      />
    </svg>
  );
}
