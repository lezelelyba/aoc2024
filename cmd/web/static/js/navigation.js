/**
 * Updates session storage with currently selected day and part to be solved
 */

function createMachine(stateMachineDefinition) {
  let currentState = stateMachineDefinition.initialState;

  return {
    getCurrentState: () => currentState,
    transition: function(event, payload) {
      const currentStateDefinition = stateMachineDefinition.states[currentState];
      const nextState = currentStateDefinition.transitions[event];

      if (!nextState) {
        return;
      }

      const prevState = currentState;
      currentState = nextState;

      if (currentStateDefinition.onExit) {
        currentStateDefinition.onExit(prevState, nextState, payload);
      }

      if (stateMachineDefinition.states[nextState].onEntry) {
        stateMachineDefinition.states[nextState].onEntry(prevState, nextState, payload);
      }

      return currentState;
    }
  };
}

const UIMachine = {
  initialState: 'idle',

  states : {
    idle: {
      transitions: {
        AUTHENTICATE: 'startAuth'
      }
    },
    startAuth: {
      transitions: {
        AUTHOK: 'authenticated',
        AUTHFAIL: 'idle'
      }
    },
    authenticated: {
      transitions: {
        SELECT: 'selected',
        TIMEOUTLOGOUT: 'idle'
      }
    },
    selected: {
      transition: {
        SUBMIT: 'submitted',
        REMOVESELECTION : 'authenticated',
        TIMEOUTLOGOUT: 'idle'
      }
    },
    submitted: {
      transition: {
        SEEN: 'authenticated',
        TIMEOUTLOGOUT: 'idle'
      }
    }
  }
}

const UIHandler = craeteMachine(UIMachine)

/**
 * Updates elements displaying currently selected day and part
 */
function updateDayPartUI() {
  const day = sessionStorage.getItem("day");
  const part = sessionStorage.getItem("part");

  const dayEl = document.getElementById("solver-day");
  const partEl = document.getElementById("solver-part");
  const solverSelectionDiv = document.getElementById("solver-selection");
  
  if (dayEl) dayEl.textContent = day || "None";
  if (partEl) partEl.textContent = part || "None";
  if (solverSelectionDiv) solverSelectionDiv.hidden = day && part ? false : true;
}

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll('a[data-day]').forEach(link => {
    link.addEventListener('click', async e => {
      e.preventDefault();
      const { day, part } = e.target.dataset;
      sessionStorage.setItem("day", day);
      sessionStorage.setItem("part", part);

      updateDayPartUI();
    });
  });
});

/**
 * Handles Submit button click
 * 
 * Check if file is selected or text area is populated, encodes input to base64
 * Calls helper to submit the input to API
 *  
 * @param {string} endpointTemplate - API endpoint template
 * @returns 
 */
async function handleSubmitClick(endpointTemplate) {
  const input = document.getElementById('fileInput');
  const file = input.files[0];

  const textInput = document.getElementById('textInput');
  const text = textInput.value;

  const resultEl = document.getElementById('result');

  let apiEndpoint
  try {
      apiEndpoint = fillTemplateFromSession(endpointTemplate)
  } catch (error) {
    resultEl.textContent = 'Error: ' + error.message;
    return
  }
  
  // check if something was filled
  // display error if not
  if (!file && text.trim() === "" ) {
    resultEl.textContent = 'Please select a file of input text first.';
    return;
  }
  
  // base64 encode the input, send request to API and display return value
  try {
    let base64;
    if (file) {
      base64 = await toBase64(file);
    } else {
      base64 = btoa(unescape(encodeURIComponent(text)));
    }
    
    const response = await sendToApi(apiEndpoint, { input: base64 })
    resultEl.textContent = JSON.stringify(response, null, 2)
  } catch(error) {
    resultEl.textContent = 'Error: ' + error.message;
  }
}

function clearFileSelection(elementId) {
  document.getElementById(elementId).value = "";
}

/**
 * Updates authentication related UI elements on the authentication state
 * 
 * if accessToken exists, session is authenticated
 * class .auth-reguired is enabled
 * elem auth-state is updated
 */
function updateAuthUI() {
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
  const loginButtonsDiv = document.getElementById("login-buttons-div");
  if (loginButtonsDiv) {
    loginButtonsDiv.hidden = accessToken ? true : false
  }

  const logoutButtonsDiv = document.getElementById("logout-buttons-div");
  if (logoutButtonsDiv) {
    logoutButtonsDiv.hidden = accessToken ? false : true
  }
}

function OAuthLogout() {
  const accessToken = sessionStorage.getItem("accessToken");
  if (accessToken) {
    sessionStorage.removeItem("accessToken");
  }

  updateAuthUI();
}