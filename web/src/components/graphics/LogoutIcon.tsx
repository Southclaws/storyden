import { SVGProps } from "react";

export function LogoutIcon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 16 16"
      width="1em"
      fill="none"
      {...props}
    >
      <path
        d="M10.5 5.83325L12.8333 7.99992L10.5 10.1666"
        stroke="#303030"
        strokeWidth="1"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M12.6667 8H7.16675"
        stroke="#303030"
        strokeWidth="1"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M10.1667 3.16675H4.50008C3.7637 3.16675 3.16675 3.7637 3.16675 4.50008V11.5001C3.16675 12.2365 3.7637 12.8334 4.50008 12.8334H10.1667"
        stroke="#303030"
        strokeWidth="1"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
