/**
 * Starts OAuth with OAuth Provider 
 * 
 * sends authorization request towards Authorization URL, adding a random state
 * 
 * @param {string} authURL - Authorization URL
 */
function startOAuth(authURL) {
  const accessToken = sessionStorage.getItem("accessToken");

  // authenticate only if token doesn't exist
  if (accessToken === null) {
    // remember current location
    sessionStorage.setItem("postAuthRedirect", window.location.pathname + window.location.search);
   
    // generate state
    const state = Math.random().toString(36).substring(2);
    sessionStorage.setItem("oauth_state", state);

    // navigate to authentication URL
    const authUrl = `${authURL}&state=${state}`;
    window.location.href = authUrl;
  }
}

/**
 * Sends OAuth code to backend to exchange it for token via API
 * 
 * code is received as query param in the redirect response from OAuth provider
 *  
 * @param {string} codeExchangeURL - Backend URL where OAuth code can be exchanged for token 
 * @param {string} provider - Name of the OAuth provider
 */
async function handleOAuthCallbackAPI(codeExchangeURL, provider) {
  // get params from callback URL
  const params = new URLSearchParams(window.location.search);
  const code = params.get("code");
  const state = params.get("state");
  const storedState = sessionStorage.getItem("oauth_state");

  // if state doesn't match abort
  if (code && state === storedState) {
    try {
      // exchange code for token
      const resp = await fetch(`${codeExchangeURL}`, { method: "POST", body: JSON.stringify({provider: provider, code: code}) });
      const data = await resp.json();

      if (data.access_token) {
          sessionStorage.setItem("accessToken", data.access_token);
      } else {
          console.warn("No access token returned from server");
      }
    } catch (err) {
      console.error("Token exchange failed", err);
    }
  }

  // remove oauth state
  sessionStorage.removeItem("oauth_state");

  // Restore original page
  const redirectTarget = sessionStorage.getItem("postAuthRedirect") || "/";
  window.history.replaceState({}, "", redirectTarget);
  sessionStorage.removeItem("postAuthRedirect");

  window.location.href = redirectTarget;
}
