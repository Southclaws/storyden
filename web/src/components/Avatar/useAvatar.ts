export function useAvatar(handle: string) {
  return { src: `/api/v1/accounts/${handle}/avatar` };
}
