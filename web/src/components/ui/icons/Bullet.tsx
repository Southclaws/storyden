import { styled } from "@/styled-system/jsx";

const Bullet = () => (
  <styled.svg
    xmlns="http://www.w3.org/2000/svg"
    width="16"
    height="16"
    viewBox="0 0 24 24"
    fill="text.subtle"
  >
    <circle cx="12.1" cy="12.1" r="2.5" />
  </styled.svg>
);

export const BulletIcon = styled(Bullet);
