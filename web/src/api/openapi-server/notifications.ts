/**
 * Generated by orval v7.2.0 🍺
 * Do not edit manually.
 * storyden
 * Storyden social API for building community driven platforms.
The Storyden API does not adhere to semantic versioning but instead applies a rolling strategy with deprecations and minimal breaking changes. This has been done mainly for a simpler development process and it may be changed to a more fixed versioning strategy in the future. Ultimately, the primary way Storyden tracks versions is dates, there are no set release tags currently.

 * OpenAPI spec version: v1.25.3-canary
 */
import type {
  NotificationListOKResponse,
  NotificationListParams,
  NotificationUpdateBody,
  NotificationUpdateOKResponse,
} from "../openapi-schema";
import { fetcher } from "../server";

/**
 * Retreive all notifications.
 */
export type notificationListResponse = {
  data: NotificationListOKResponse;
  status: number;
};

export const getNotificationListUrl = (params?: NotificationListParams) => {
  const normalizedParams = new URLSearchParams();

  Object.entries(params || {}).forEach(([key, value]) => {
    if (value !== undefined) {
      normalizedParams.append(key, value === null ? "null" : value.toString());
    }
  });

  return normalizedParams.size
    ? `/notifications?${normalizedParams.toString()}`
    : `/notifications`;
};

export const notificationList = async (
  params?: NotificationListParams,
  options?: RequestInit,
): Promise<notificationListResponse> => {
  return fetcher<Promise<notificationListResponse>>(
    getNotificationListUrl(params),
    {
      ...options,
      method: "GET",
    },
  );
};

/**
 * Change the read status for a notification.
 */
export type notificationUpdateResponse = {
  data: NotificationUpdateOKResponse;
  status: number;
};

export const getNotificationUpdateUrl = (notificationId: string) => {
  return `/notifications/${notificationId}`;
};

export const notificationUpdate = async (
  notificationId: string,
  notificationUpdateBody: NotificationUpdateBody,
  options?: RequestInit,
): Promise<notificationUpdateResponse> => {
  return fetcher<Promise<notificationUpdateResponse>>(
    getNotificationUpdateUrl(notificationId),
    {
      ...options,
      method: "PATCH",
      headers: { "Content-Type": "application/json", ...options?.headers },
      body: JSON.stringify(notificationUpdateBody),
    },
  );
};
