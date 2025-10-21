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
  // beginning state, cannot be revisited
  states : {
    start: {
      transitions: {
        START: 'idle'
      },
    },
  // resets all elements to default state
    idle: {
      transitions: {
        AUTHENTICATE: 'startAuth',
        SKIPAUTH: 'authenticated'
      },
      onEntry: function(prevState, thisState, payload) {

        // auth
        const authEls = document.getElementsByClassName("auth");
        Array.from(authEls).forEach(el => {
          el.classList.add("active");
          el.classList.remove("done");
        });
        
        const statusEl = document.getElementById("auth-status");
        if(statusEl) {
          statusEl.textContent = "Not Authenticated";
          statusEl.classList.remove("authenticated");
        }
        
        const loginButtonsDiv = document.getElementById("login-buttons-div");
        if(loginButtonsDiv) loginButtonsDiv.hidden = false;

        const logoutButtonsDiv = document.getElementById("logout-buttons-div");
        if(logoutButtonsDiv) logoutButtonsDiv.hidden = true;

        // selection
        const selectionEls = document.getElementsByClassName("selection");
        Array.from(selectionEls).forEach(el => {
          el.classList.add("disabled");
          el.classList.remove("active");
          el.classList.remove("done");
        });
        
        sessionStorage.removeItem("day");
        sessionStorage.removeItem("part");

        const solverSelectionDiv = document.getElementById("solver-selection");
        const dayEl = document.getElementById("solver-day");
        const partEl = document.getElementById("solver-part");
        
        if (dayEl) dayEl.textContent = "None";
        if (partEl) partEl.textContent = "None";
        if (solverSelectionDiv) solverSelectionDiv.hidden = true;

        // submission
        const submissionEls= document.getElementsByClassName("submission");
        Array.from(submissionEls).forEach(el => {
          el.classList.add("disabled");
          el.classList.remove("active");
          el.classList.remove("done");
        });

        document.getElementById("fileInput").value = "";
        document.getElementById("textInput").value = "";
        const fileNameEl = document.getElementById('fileName')
        fileNameEl.textContent = fileNameEl.dataset.default;

        // result
        const resultEls = document.getElementsByClassName("result");
        Array.from(resultEls).forEach(el => {
          el.classList.add("disabled");
          el.classList.remove("active");
          el.classList.remove("done");
        });

        const seenBtn = document.getElementById("seenBtn").hidden = true;
        const resultEl = document.getElementById('result').textContent = "";
      }
    },
    // starts authentication - currently only a transition state
    startAuth: {
      transitions: {
        AUTHOK: 'authenticated',
        AUTHFAIL: 'idle'
      }
    },
    // user is authenticated
    authenticated: {
      transitions: {
        SELECT: 'selected',
        TIMEOUTLOGOUT: 'idle'
      },
      onEntry: function(prevState, thisState, payload) {
        // auth

        // display authenticated message
        const statusEl = document.getElementById("auth-status");
        if(statusEl) {
          statusEl.textContent = "Authenticated";
          statusEl.classList.add("authenticated");
        }
       
        // hide login buttons
        const loginButtonsDiv = document.getElementById("login-buttons-div");
        if(loginButtonsDiv) loginButtonsDiv.hidden = true;

        // show logout button
        const logoutButtonsDiv = document.getElementById("logout-buttons-div");
        if(logoutButtonsDiv) logoutButtonsDiv.hidden = false;

        // mark authentication sections as done
        const authEls = document.getElementsByClassName("auth");
        Array.from(authEls).forEach(el => {
          el.classList.remove("active");
          el.classList.add("done");
        });

        // activate selection section
        const selectionEls = document.getElementsByClassName("selection");
        Array.from(selectionEls).forEach(el => {
          el.classList.remove("disabled");
          el.classList.add("active");
        });
      },
    },
    // file is selected
    selected: {
      transitions: {
        SUBMIT: 'submitted',
        SELECT: 'selected',
        TIMEOUTLOGOUT: 'idle'
      },
      onEntry: function(prevState, thisState, payload) {
        // get selected day and part from storage (reselect)
        let day = sessionStorage.getItem("day");
        let part = sessionStorage.getItem("part");

        // get day and part from payload
        if (payload !== undefined && ("day" in payload && "part" in payload)) {
          sessionStorage.setItem("day", payload.day);
          sessionStorage.setItem("part", payload.part);

          day = payload.day;
          part = payload.part;
        }
       
        // update selection display
        const dayEl = document.getElementById("solver-day");
        const partEl = document.getElementById("solver-part");
        const solverSelectionDiv = document.getElementById("solver-selection");
        
        if (dayEl) dayEl.textContent = day || "None";
        if (partEl) partEl.textContent = part || "None";
        if (solverSelectionDiv) solverSelectionDiv.hidden = day && part ? false : true;
 
        // mark selection as done
        const selectionEls = document.getElementsByClassName("selection");
        Array.from(selectionEls).forEach(el => {
          el.classList.remove("active");
          el.classList.add("done");
        });

        // activate submit section
        const submissionEls = document.getElementsByClassName("submission");
        Array.from(submissionEls).forEach(el => {
          el.classList.remove("disabled");
          el.classList.add("active");
        });
      },
    },
    // show result
    submitted: {
      transitions: {
        SEEN: 'idle',
        TIMEOUTLOGOUT: 'idle',
        SUBMITERROR: 'selected'
      },
      onEntry: function(prevState, thisState, payload) {
        // disable selection until result is acknowledged
        const selectionEls = document.getElementsByClassName("selection");
        Array.from(selectionEls).forEach(el => {
          el.classList.remove("active");
          el.classList.add("done");
          el.classList.add("disabled");
        });

        // disable submit section until result is acknowledged
        const submissionEls = document.getElementsByClassName("submission");
        Array.from(submissionEls).forEach(el => {
          el.classList.remove("active");
          el.classList.add("done");
          el.classList.add("disabled");
        });

        // activate result section
        const resultEls = document.getElementsByClassName("result");
        Array.from(resultEls).forEach(el => {
          el.classList.remove("disabled");
          el.classList.add("active");
        });

        // show seen button
        const seenBtn = document.getElementById("seenBtn");
        seenBtn.hidden = false;
      },
      onExit: function(thisState, nextState, payload) {
        // local submit error (wrong file, no input)
        // init would reset everything
        // 
        if (nextState === "selected") {

          // hide seen button
          const seenBtn = document.getElementById("seenBtn");
          seenBtn.hidden = true ;

          // clan up result
          const resultEl = document.getElementById('result');
          resultEl.textContent = "";

          // disable result
          const resultEls = document.getElementsByClassName("result");
          Array.from(resultEls).forEach(el => {
            el.classList.remove("active");
            el.classList.add("disabled");
          });

          // reenable selection, keep the day/part
          const selectionEls = document.getElementsByClassName("selection");
          Array.from(selectionEls).forEach(el => {
            el.classList.remove("disabled");
            el.classList.add("done");
            el.classList.add("active");
          });
          
        }
      }
    }
  }
}

// create Machine
const UIHandler = createMachine(UIMachine);

// register events

// on load check authentication state, as we can be redirected from callback
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

// logout button
// might not exist if auth is disabled
const logoutBtn = document.getElementById("logoutBtn")
if(logoutBtn) {
  logoutBtn.addEventListener("click", () => {
    sessionStorage.removeItem("accessToken");
    UIHandler.transition('TIMEOUTLOGOUT');
  });
}

// day/part selection table
document.querySelectorAll('a[data-day]').forEach(link => {
  link.addEventListener('click', e => {
    e.preventDefault();
    const { day, part } = e.target.dataset;

    UIHandler.transition('SELECT', {day: day, part: part});
  });
});

// submit button
// has to be async
document.getElementById("submitBtn").addEventListener("click", async () => {
  // reset seenBtn
  // clicking goes to init and resets display
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
  // if submit threw error -> modify seen button to allow fix and resubmit
  // if submit was ok -> keep original to go to init state
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

// file input change
// modifies selected file display
document.getElementById("fileInput").addEventListener("change", e => {
  fileNameEl = document.getElementById('fileName')
  const fileName = e.target.files[0]?.name || fileNameEl.dataset.default;
  fileNameEl.textContent = fileName;
});

// input clear button
// clears file and text inputs, sets text input to default
document.getElementById("inputClear").addEventListener("click", () => {
  document.getElementById("fileInput").value = "";
  document.getElementById("textInput").value = "";
  fileNameEl = document.getElementById('fileName')
  fileNameEl.textContent = fileNameEl.dataset.default;
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