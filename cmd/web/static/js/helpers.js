/**
 * Fills template string from session storage
 * 
 * replace {key} placeholders within template with value found under "key" in session storage
 * @param {string} template - template
 * @returns {string} Filled template
 * @throws {Error} If some key is missing from session storage
 */
function fillTemplateFromSession(template) {
  return template.replace(/{(\w+)}/g, (match, key) => {
    const value = sessionStorage.getItem(key);
    if (value === null) {
      throw new Error(`Missing session value for: ${key}`);
    }
    return value;
  });
}

/**
 * Parses JWT token and returns claims
 * @param {string} token - JWT Token
 * @returns {object} Claims
 */
function parseJwt(token) {
  const payload = token.split('.')[1];
  return JSON.parse(atob(payload));
}


function startAuthTimer(fsm, elId, exp) {
  let timerId;

  function updateTimer() {
    const now = Math.floor(Date.now() / 1000);
    let remaining = exp - now;
    if (remaining <= 0) {
      remaining = 0;
      if (timerId) {
        clearInterval(timerId);
        fsm.transition('TIMEOUTLOGOUT', {timeout: true});
      }
    }

    const minutes = Math.floor(remaining / 60);
    const seconds = remaining % 60;
    document.getElementById(elId).textContent =
      `${minutes}:${seconds.toString().padStart(2, "0")}`;
  }

  updateTimer();
  timerId = setInterval(updateTimer, 1000);
}