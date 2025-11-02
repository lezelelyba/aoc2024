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
  return timerId
}