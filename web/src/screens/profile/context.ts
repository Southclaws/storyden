import { createContext, useContext } from "react";

export type ProfileContextShape = {
  isSelf: boolean;
  isEditing: boolean;
};

export const ProfileContext = createContext<ProfileContextShape>({
  isSelf: false,
  isEditing: false,
});

export function useProfileContext() {
  return useContext(ProfileContext);
}
