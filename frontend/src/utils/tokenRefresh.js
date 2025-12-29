import { refreshToken } from "../api";
import { getAccessToken, setAccessToken } from "./tokenStorage";

let refreshTimer = null;
const REFRESH_INTERVAL = 20 * 60 * 1000; // Refresh every 20 minutes (access tokens last 30 minutes)

export function startTokenRefresh() {
  // Clear any existing timer
  if (refreshTimer) {
    clearInterval(refreshTimer);
  }

  // Only start refresh if we have an access token
  if (!getAccessToken()) {
    return;
  }

  // Refresh token immediately on start
  refreshToken().catch((err) => {
    console.error("Initial token refresh failed:", err);
  });

  // Set up periodic refresh
  refreshTimer = setInterval(() => {
    // Only refresh if we still have an access token
    if (getAccessToken()) {
      refreshToken().catch((err) => {
        console.error("Token refresh failed:", err);
        // If refresh fails, user will need to login again
        // Could redirect to login page here if needed
      });
    } else {
      // No access token, stop refreshing
      stopTokenRefresh();
    }
  }, REFRESH_INTERVAL);
}

export function stopTokenRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer);
    refreshTimer = null;
  }
}

// Also refresh token before it expires (refresh when there's 1 day left)
export function setupTokenRefresh() {
  // Only refresh if we have an access token
  if (getAccessToken()) {
    // Refresh token on page load
    refreshToken().catch((err) => {
      console.error("Token refresh on page load failed:", err);
    });

    // Set up periodic refresh
    startTokenRefresh();
  }
}







