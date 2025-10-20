/**
 * State Machine 
 * @param {object} stateMachineDefinition - Definitioin of the state machine
 * @returns {object} - state machine
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

/**
 * State Machine definition
 */
const UIMachine = {
  initialState: 'start',
  states : {
    start: {
      transitions: {
        START: 'idle'
      },
    },
    idle: {
      transitions: {
        AUTHENTICATE: 'startAuth',
        SKIPAUTH: 'authenticated'
      },
      onEntry: function(prevState, thisState, payload) {

        const authEls = document.getElementsByClassName("auth");

        Array.from(authEls).forEach(el => {
          el.classList.add("active");
          el.classList.remove("done");
        });
        
        const statusEl = document.getElementById("auth-status");
        if(statusEl) statusEl.textContent = "Not Authenticated";
        
        const loginButtonsDiv = document.getElementById("login-buttons-div");
        if(loginButtonsDiv) loginButtonsDiv.hidden = false;

        const logoutButtonsDiv = document.getElementById("logout-buttons-div");
        if(logoutButtonsDiv) logoutButtonsDiv.hidden = true;

        const selectionEls = document.getElementsByClassName("selection");
        Array.from(selectionEls).forEach(el => {
          el.classList.add("disabled");
          el.classList.remove("active");
          el.classList.remove("done");
          el.classList.remove("donedisabled");
        });
        
        sessionStorage.removeItem("day");
        sessionStorage.removeItem("part");

        const solverSelectionDiv = document.getElementById("solver-selection");
        const dayEl = document.getElementById("solver-day");
        const partEl = document.getElementById("solver-part");
        
        if (dayEl) dayEl.textContent = "None";
        if (partEl) partEl.textContent = "None";
        if (solverSelectionDiv) solverSelectionDiv.hidden = true;

        const submissionEls= document.getElementsByClassName("submission");
        Array.from(submissionEls).forEach(el => {
          el.classList.add("disabled");
          el.classList.remove("active");
          el.classList.remove("done");
          el.classList.remove("donedisabled");
        });

        document.getElementById("fileInput").value = "";
        document.getElementById("textInput").value = "";

        const resultEls = document.getElementsByClassName("result");
        Array.from(resultEls).forEach(el => {
          el.classList.add("disabled");
          el.classList.remove("active");
          el.classList.remove("done");
          el.classList.remove("donedisabled");
        });

        const seenBtn = document.getElementById("seenBtn");
        seenBtn.hidden = true;
        
        const resultEl = document.getElementById('result');
        resultEl.textContent = "";
        
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
      },
      onEntry: function(prevState, thisState, payload) {
        const statusEl = document.getElementById("auth-status");
        if(statusEl) statusEl.textContent = "Authenticated";
        
        const loginButtonsDiv = document.getElementById("login-buttons-div");
        if(loginButtonsDiv) loginButtonsDiv.hidden = true;

        const logoutButtonsDiv = document.getElementById("logout-buttons-div");
        if(logoutButtonsDiv) logoutButtonsDiv.hidden = false;

        const authEls = document.getElementsByClassName("auth");
        Array.from(authEls).forEach(el => {
          el.classList.remove("active");
          el.classList.add("done");
        });

        const selectionEls = document.getElementsByClassName("selection");
        Array.from(selectionEls).forEach(el => {
          el.classList.remove("disabled");
          el.classList.add("active");
        });
      },
    },
    selected: {
      transitions: {
        SUBMIT: 'submitted',
        SELECT: 'selected',
        TIMEOUTLOGOUT: 'idle'
      },
      onEntry: function(prevState, thisState, payload) {
        let day = sessionStorage.getItem("day");
        let part = sessionStorage.getItem("part");

        if (payload !== undefined && ("day" in payload && "part" in payload)) {
          sessionStorage.setItem("day", payload.day);
          sessionStorage.setItem("part", payload.part);

          day = payload.day;
          part = payload.part;
        }
        
        const dayEl = document.getElementById("solver-day");
        const partEl = document.getElementById("solver-part");
        const solverSelectionDiv = document.getElementById("solver-selection");
        
        if (dayEl) dayEl.textContent = day || "None";
        if (partEl) partEl.textContent = part || "None";
        if (solverSelectionDiv) solverSelectionDiv.hidden = day && part ? false : true;
  
        const selectionEls = document.getElementsByClassName("selection");
        Array.from(selectionEls).forEach(el => {
          el.classList.remove("active");
          el.classList.remove("donedisabled");
          el.classList.add("done");
        });

        const submissionEls = document.getElementsByClassName("submission");
        Array.from(submissionEls).forEach(el => {
          el.classList.remove("disabled");
          el.classList.remove("donedisabled");
          el.classList.add("active");
        });
      },
    },
    submitted: {
      transitions: {
        SEEN: 'idle',
        TIMEOUTLOGOUT: 'idle',
        SUBMITERROR: 'selected'
      },
      onEntry: function(prevState, thisState, payload) {
        const selectionEls = document.getElementsByClassName("selection");
        Array.from(selectionEls).forEach(el => {
          el.classList.remove("active");
          el.classList.remove("done");
          el.classList.add("donedisabled");
        });

        const submissionEls = document.getElementsByClassName("submission");
        Array.from(submissionEls).forEach(el => {
          el.classList.remove("active");
          el.classList.add("donedisabled");
        });

        const resultEls = document.getElementsByClassName("result");
        Array.from(resultEls).forEach(el => {
          el.classList.remove("disabled");
          el.classList.add("active");
        });

        const seenBtn = document.getElementById("seenBtn");
        seenBtn.hidden = false;
      },
      onExit: function(thisState, nextState, payload) {
        if (nextState === "selected") {
          const seenBtn = document.getElementById("seenBtn");
          seenBtn.hidden = true ;

          const resultEls = document.getElementsByClassName("result");
          Array.from(resultEls).forEach(el => {
            el.classList.add("disabled");
            el.classList.remove("active");
          });
          
          const resultEl = document.getElementById('result');
          resultEl.textContent = "";
        }
      }
    }
  }
}

// create Machine
const UIHandler = createMachine(UIMachine);

// register events
document.addEventListener("DOMContentLoaded", function() {
    console.log("DOM fully loaded and parsed");

    UIHandler.transition('START');

    const configEl = document.getElementById("auth-enabled");
    const authEnabled = configEl.dataset.enabled;
    
    if (authEnabled === "true") {
      const accessToken = sessionStorage.getItem("accessToken");

      // authenticated
      if (accessToken) {
        UIHandler.transition('AUTHENTICATE');
        UIHandler.transition('AUTHOK');
      }
    } else {
        accessToken = sessionStorage.removeItem("accessToken");
        UIHandler.transition('SKIPAUTH');
    }
});

// might not exist if auth is disabled
const logoutBtn = document.getElementById("logoutBtn")
if(logoutBtn) {
  logoutBtn.addEventListener("click", () => {
    sessionStorage.removeItem("accessToken");
    UIHandler.transition('TIMEOUTLOGOUT');
  });
}

document.querySelectorAll('a[data-day]').forEach(link => {
  link.addEventListener('click', async e => {
    e.preventDefault();
    const { day, part } = e.target.dataset;

    UIHandler.transition('SELECT', {day: day, part: part});
  });
});

document.getElementById("submitBtn").addEventListener("click", async () => {
  // setup seenBtn

  // remove old listener
  const seenBtn = document.getElementById("seenBtn");
  seenBtn.replaceWith(seenBtn.cloneNode(true));

  // put in a new one
  document.getElementById("seenBtn").addEventListener("click", () => {
    UIHandler.transition('SEEN');
    
    const configEl = document.getElementById("auth-enabled");
    const authEnabled = configEl.dataset.enabled;
    
    if (authEnabled === "true") {
      const accessToken = sessionStorage.getItem("accessToken");

      // authenticated
      if (accessToken) {
        UIHandler.transition('AUTHENTICATE');
        UIHandler.transition('AUTHOK');
      }
    } else {
        accessToken = sessionStorage.removeItem("accessToken");
        UIHandler.transition('SKIPAUTH');
    }
  });
 
  // submit
  UIHandler.transition('SUBMIT');
  // check if submit failed
  try {
    await handleSubmitClick("/api/solvers/{day}/{part}");
  } catch(error) {
    // setup seenBtn to return to submit state
    // remove old listener
    const seenBtn = document.getElementById("seenBtn");
    seenBtn.replaceWith(seenBtn.cloneNode(true));

    document.getElementById("seenBtn").addEventListener("click", () => {
      UIHandler.transition('SUBMITERROR');
    });
  }
});

// document.getElementById("textInput").addEventListener("input", e => {
//   fileInput = document.getElementById("fileInput")
//   submitBtn = document.getElementById("submitBtn")
// 
//   submitBtn.disabled = (e.target.value.trim() === "" ) && (fileInput.value === "");
// });
// 
// document.getElementById("fileInput").addEventListener("change", e => {
//   textInput = document.getElementById("textInput")
//   submitBtn = document.getElementById("submitBtn")
// 
//   submitBtn.disabled = (textInput.value.trim() === "" ) && (e.target.value === "");
// });

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
    resultEl.textContent = 'Local Error: ' + error.message;
    throw(error);
  }
  
  // check if something was filled
  // display error if not
  if (!file && text.trim() === "" ) {
    resultEl.textContent = 'Please select a file of input text first.';
    throw Error("no input");
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
    throw(error);
  }
}

// function clearFileSelection(elementId) {
//   document.getElementById(elementId).value = "";
// }

function clearSelection() {
  document.getElementById("fileInput").value = "";
  document.getElementById("textInput").value = "";
}