import { Box } from "@/styled-system/jsx";

import styles from "./navigation.module.css";

import { getServerSidebarState } from "./NavigationPane/server";

type Props = {
  children: React.ReactNode;
};

export async function NavigationContainer({ children }: Props) {
  const showLeftBar = await getServerSidebarState();

  return (
    <Box
      id="navigation__container"
      className={styles["navigation__container"]}
      w="full"
      data-leftbar-shown={showLeftBar}
    >
      {children}
    </Box>
  );
}
