// Use environment variable, or empty string for same-origin (production), or localhost for dev
const API = import.meta.env.VITE_API_BASE || (import.meta.env.PROD ? "" : "http://localhost:8080");

export async function getMe() {
  const res = await fetch(`${API}/me`, { credentials: "include" });
  if (!res.ok) return null;
  return res.json();
}
