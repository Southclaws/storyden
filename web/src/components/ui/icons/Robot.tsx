import { styled } from "@/styled-system/jsx";

const Clanker = (props: any) => (
  <svg
    width="24"
    height="24"
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    {...props}
  >
    <path
      d="M11.9995 8.1411C11.9995 8.1411 13.2761 3.35972 11.9995 4.14111C10.9001 4.81411 8.96312 4.92249 8.0293 3.22256"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M18 8.14111H6C4.89543 8.14111 4 9.03654 4 10.1411V18.1411C4 19.2457 4.89543 20.1411 6 20.1411H18C19.1046 20.1411 20 19.2457 20 18.1411V10.1411C20 9.03654 19.1046 8.14111 18 8.14111Z"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M2 14.1411C3.11764 14.4721 4 14.1411 4 14.1411"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M20 14.1411C20 14.1411 21.0076 13.8844 22 14.1411"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M15 12.1411V14.1411"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M9 12.1411V14.1411"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const RobotIcon = styled(Clanker);
