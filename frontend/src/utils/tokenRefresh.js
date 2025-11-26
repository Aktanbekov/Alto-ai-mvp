import { refreshToken } from "../api";

let refreshTimer = null;
const REFRESH_INTERVAL = 6 * 60 * 60 * 1000; // Refresh every 6 hours (tokens last 7 days)

export function startTokenRefresh() {
  // Clear any existing timer
  if (refreshTimer) {
    clearInterval(refreshTimer);
  }

  // Refresh token immediately on start
  refreshToken().catch((err) => {
    console.error("Initial token refresh failed:", err);
  });

  // Set up periodic refresh
  refreshTimer = setInterval(() => {
    refreshToken().catch((err) => {
      console.error("Token refresh failed:", err);
      // If refresh fails, user will need to login again
      // Could redirect to login page here if needed
    });
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
  // Refresh token on page load
  refreshToken().catch((err) => {
    console.error("Token refresh on page load failed:", err);
  });

  // Set up periodic refresh
  startTokenRefresh();
}

