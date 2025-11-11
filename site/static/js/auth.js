/**
 * Starts OAuth with OAuth Provider 
 * 
 * sends authorization request towards Authorization URL, adding a random state
 * 
 * @param {string} authURL - Authorization URL
 */
function startOAuth(authURL) {
  // do not try to authenticate if no backend is available
  // call back exchanges code with backend
  let backend = null;
  
  if (typeof getBackend === "function") {
    backend = getBackend();
  }

  if (!backend) {
    console.log("No backend available");
    return;
  }

  // authenticate only if token doesn't exist
  if (getAccessToken() === null) {

    // set backend for call back
    sessionStorage.setItem("baseApiUrl", backend.baseApiUrl);
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
  const baseApiUrl = sessionStorage.getItem("baseApiUrl");

  // if state doesn't match abort
  if (code && state === storedState) {
    try {
      // exchange code for token
      const resp = await fetch(`${baseApiUrl}${codeExchangeURL}`, { method: "POST", body: JSON.stringify({provider: provider, code: code}) });
      const data = await resp.json();

      if (data.access_token) {
          setAccessToken(data.access_token);
      } else {
          dataStr = JSON.stringify(data);
          console.warn(`No access token returned from server: ${dataStr}`);
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

function authRequired() {
  const configEl = document.getElementById("auth-enabled");

  if (configEl && configEl.dataset.enabled == "true") {
    return true
  }
 
  return false
}

function getAccessToken() {
  const accessToken = sessionStorage.getItem("accessToken");

  return accessToken
}

function setAccessToken(token) {
  sessionStorage.setItem("accessToken", token);
}

function removeAccessToken() {
  sessionStorage.removeItem("accessToken");
}

/*
 * Parses JWT token and returns claims
 * @param {string} token - JWT Token
 * @returns {object} Claims
 */
function parseJwt(token) {
  const payload = token.split('.')[1];
  return JSON.parse(atob(payload));
}


/**
 * Starts Authentication timer
 * Transitions via TIMEOUTLOGOUT after timer expires
 *  
 * @param {number} exp - expiration time of authentication
 * @param {string} elId - element where to write the remaining time
 * @param {function} expFunc - expiration function
 * @returns {object} Timer
 */
function startAuthTimer(exp, elId, expFunc) {
  let timerId;

  // create a function which is called periodically
  function updateTimer() {
    const now = Math.floor(Date.now() / 1000);
    let remaining = exp - now;

    // when expired
    if (remaining <= 0) {
      remaining = 0;
      // if timer var exists => timer is running
      if (timerId) {
        // clear
        clearInterval(timerId);
        // call the expiration function
        expFunc()
      }
    }

    // display the remaining time in the element
    const minutes = Math.floor(remaining / 60);
    const seconds = remaining % 60;
    const el = document.getElementById(elId);
    if (el) {
      el.textContent = `${minutes}:${seconds.toString().padStart(2, "0")}`;
    }
  }

  // 1st run to display the remaining time immediately
  updateTimer();
  // start timer with 1s update interval
  timerId = setInterval(updateTimer, 1000);

  // return timer
  return timerId;
}
