function updateUI() {
  const accessToken = sessionStorage.getItem("accessToken");
  const authEls = document.querySelectorAll(".auth-required");
  authEls.forEach(el => {
    if (accessToken) {
      el.disabled = false;
    } else {
      el.disabled = true;
    }
  });

  const statusEl = document.getElementById("auth-status");
  if (statusEl) {
    statusEl.textContent = accessToken ? "Authenticated" : "Not authenticated";
  }
}

function startOAuth(authURL, clientId, redirectUri) {
  const accessToken = sessionStorage.getItem("accessToken");
  if (accessToken === null) {
    sessionStorage.setItem("postAuthRedirect", window.location.pathname + window.location.search);
    
    const state = Math.random().toString(36).substring(2);
    sessionStorage.setItem("oauth_state", state);

    // const authUrl = `https://github.com/login/oauth/authorize?` +
    const authUrl = `${authURL}?` +
                    `client_id=${clientId}&redirect_uri=${encodeURIComponent(redirectUri)}` +
                    `&scope=read:user,user:email&state=${state}`;
    window.location.href = authUrl;
  }
}

async function handleOAuthCallback(endpoint) {
  const params = new URLSearchParams(window.location.search);
  const code = params.get("code");
  const state = params.get("state");
  const storedState = sessionStorage.getItem("oauth_state");

  if (code && state === storedState) {
    try {
      const resp = await fetch(`${endpoint}?code=${code}`, { method: "POST" });
      const data = await resp.json();
      sessionStorage.setItem("accessToken", data.access_token);

      sessionStorage.removeItem("oauth_state");

      // Restore original page after login
      const redirectTarget = sessionStorage.getItem("postAuthRedirect") || "/";
      window.history.replaceState({}, "", redirectTarget);
      sessionStorage.removeItem("postAuthRedirect");

      window.location.href = redirectTarget;
    } catch (err) {
      console.error("Token exchange failed", err);
    }
  }
}