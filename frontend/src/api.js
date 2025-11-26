// Use environment variable, or empty string for same-origin (production), or localhost for dev
const API = import.meta.env.VITE_API_BASE || (import.meta.env.PROD ? "" : "http://localhost:8080");

export async function getMe() {
  const res = await fetch(`${API}/me`, { credentials: "include" });
  if (!res.ok) return null;
  return res.json();
}

export async function sendChatMessage(messages) {
  const res = await fetch(`${API}/api/v1/chat`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ messages }),
  });
  
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: "Failed to get response" }));
    throw new Error(error.error || "Failed to send message");
  }
  
  const data = await res.json();
  return data.data.content;
}

export async function login(email, password) {
  const res = await fetch(`${API}/api/v1/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ email, password }),
  });
  
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: "Login failed" }));
    throw new Error(error.error || "Login failed");
  }
  
  return res.json();
}

export async function logout() {
  const res = await fetch(`${API}/api/v1/auth/logout`, {
    method: "POST",
    credentials: "include",
  });
  
  if (!res.ok) {
    throw new Error("Logout failed");
  }
}

export async function refreshToken() {
  const res = await fetch(`${API}/api/v1/auth/refresh`, {
    method: "POST",
    credentials: "include",
  });
  
  if (!res.ok) {
    throw new Error("Token refresh failed");
  }
  
  return res.json();
}

export async function register(email, name, password) {
  try {
    const res = await fetch(`${API}/api/v1/auth/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ email, name, password }),
    });
    
    // Check if response is HTML (wrong route)
    const contentType = res.headers.get("content-type");
    if (contentType && contentType.includes("text/html")) {
      throw new Error("Server returned HTML instead of JSON. The API endpoint may not be registered. Please restart the server.");
    }
    
    if (!res.ok) {
      let error;
      try {
        error = await res.json();
      } catch (e) {
        throw new Error(`Registration failed: ${res.status} ${res.statusText}`);
      }
      
      // Handle validation errors
      if (error.error === "validation_error" && error.details) {
        const details = error.details;
        const messages = [];
        
        if (details.email) messages.push(`Email: ${getValidationMessage(details.email)}`);
        if (details.name) messages.push(`Name: ${getValidationMessage(details.name)}`);
        if (details.password) messages.push(`Password: ${getValidationMessage(details.password)}`);
        
        throw new Error(messages.length > 0 ? messages.join(". ") : "Validation failed");
      }
      
      throw new Error(error.error || "Registration failed");
    }
    
    return res.json();
  } catch (err) {
    // Handle network errors
    if (err instanceof TypeError && err.message.includes("fetch")) {
      throw new Error("Unable to connect to server. Please check if the server is running.");
    }
    throw err;
  }
}

function getValidationMessage(tag) {
  const messages = {
    required: "is required",
    email: "must be a valid email address",
    min: "is too short",
    max: "is too long",
    len: "has incorrect length",
  };
  return messages[tag] || `failed validation (${tag})`;
}

export async function verifyEmail(email, code) {
  const res = await fetch(`${API}/api/v1/auth/verify-email`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ email, code }),
  });
  
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: "Verification failed" }));
    throw new Error(error.error || "Verification failed");
  }
  
  return res.json();
}

export async function resendVerificationCode(email) {
  const res = await fetch(`${API}/api/v1/auth/resend-verification`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ email }),
  });
  
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: "Request failed" }));
    throw new Error(error.error || "Request failed");
  }
  
  return res.json();
}

export async function forgotPassword(email) {
  const res = await fetch(`${API}/api/v1/auth/forgot-password`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ email }),
  });
  
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: "Request failed" }));
    throw new Error(error.error || "Request failed");
  }
  
  return res.json();
}

export async function resetPassword(email, code, password) {
  const res = await fetch(`${API}/api/v1/auth/reset-password`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify({ email, code, password }),
  });
  
  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: "Password reset failed" }));
    throw new Error(error.error || "Password reset failed");
  }
  
  return res.json();
}
