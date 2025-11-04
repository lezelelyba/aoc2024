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
        // remove token if user logged out or token timed out
        if (payload && payload.timeout) {
          sessionStorage.removeItem("accessToken");
        }

        if (payload && payload.logout) {
          sessionStorage.removeItem("accessToken");
        }
        
        // enable body if timeout happend during waiting for result
        document.body.classList.remove("busy");

        // auth
        const authEls = document.getElementsByClassName("auth");
        Array.from(authEls).forEach(el => {
          el.hidden = true;
          el.classList.add("active");
          el.classList.remove("disabled");
          el.classList.remove("done");
        });
        
        const statusEl = document.getElementById("auth-status");
        if(statusEl) {
          statusEl.textContent = "Not Authenticated";
          statusEl.classList.remove("authenticated");
        }
        
        const timerEl = document.getElementById("auth-timer");
        if(timerEl) {
          timerEl.hidden = true
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

        document.getElementById("seenBtn").hidden = true;
        document.getElementById("resubmitBtn").hidden = true;
        document.getElementById('result').textContent = "";

        // check backend status
        document.body.classList.add("busy");
        sendToApi("GET", "/api/public/info")
        // update backend information
        .then(info => {
          const configEl = document.getElementById("auth-enabled");
          const backendEl= document.getElementById("backend");
          if (info.authentication == "oauth") {
            configEl.dataset.enabled = "true";
            const authEls = document.getElementsByClassName("auth");
            Array.from(authEls).forEach(el => {
              el.hidden = false;
            });
          } else {
            configEl.dataset.enabled = "false";
          }
          backendEl.classList.remove("backend-not-available");
          backendEl.classList.add("backend-available");
        })
        .catch(err => {
          const configEl = document.getElementById("auth-enabled");
          configEl.dataset.enabled = "false";
          
          const backendEl= document.getElementById("backend");
          backendEl.classList.add("backend-not-available");
          backendEl.classList.remove("backend-available");
        })
        .then(() => {
          // if backend is not available do not progress further
          const backendEl= document.getElementById("backend");
          if (backendEl.classList.contains("backend-not-available")) {
            return
          }
          const configEl = document.getElementById("auth-enabled");
          const authEnabled = configEl.dataset.enabled;
         
          // if authentication is required
          if (authEnabled === "true") {
            const accessToken = sessionStorage.getItem("accessToken");

            // if we already have token => authenticate
            if (accessToken) {
              setTimeout(() => UIHandler.transition('AUTHENTICATE'), 0);
            }
          // if not skip auth
          } else {
              accessToken = sessionStorage.removeItem("accessToken");
              setTimeout(() => UIHandler.transition('SKIPAUTH'), 0);
          }
        })
        .finally(() => {
          document.body.classList.remove("busy");
        });
      }
    },
    // starts authentication - currently only a transition state
    startAuth: {
      transitions: {
        AUTHOK: 'authenticated',
        AUTHFAIL: 'idle'
      },
      onEntry: function(prevState, thisState, payload) {
        // if we already have token => authenticate
        const accessToken = sessionStorage.getItem("accessToken");
        if (accessToken) {
          setTimeout(() => UIHandler.transition('AUTHOK'), 0);
        }
      }
    },
    // user is authenticated
    authenticated: {
      transitions: {
        SELECT: 'selected',
        TIMEOUTLOGOUT: 'idle'
      },
      onEntry: function(prevState, thisState, payload) {

        // display authenticated message
        const statusEl = document.getElementById("auth-status");
        if(statusEl) {
          statusEl.textContent = "Authenticated";
          statusEl.classList.add("authenticated");
        }

        // display authentication timer
        const timerEl = document.getElementById("auth-timer");

        if(timerEl && sessionStorage.getItem("accessToken")) {
          timerEl.hidden = false;
          const tokenExpiration = parseJwt(sessionStorage.getItem("accessToken")).exp;

          function timeout() {
            UIHandler.transition("TIMEOUTLOGOUT", {timeout: true});
          }

          startAuthTimer(tokenExpiration, "auth-timer", timeout);
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

        // display selection table if not already displayed
        if (!document.getElementById("table-container").firstElementChild) {
          sendToApi("GET", "/api/solvers", {} )
            .catch(err => {})
            .then(d => createTable("table-container", d, solverListingHeaderFunc, solverListingRowFunc))
            .then(() => {
              // register listeners for the table
              // day/part selection table
              document.querySelectorAll('a[data-day]').forEach(link => {
                link.addEventListener('click', e => {
                  e.preventDefault();
                  const { day, part } = e.target.dataset;

                  UIHandler.transition('SELECT', {day: day, part: part});
                });
              });
          });
        }

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
    // transition state to show different elements based on the results from the api
    submitted: {
      transitions: {
        SHOWRESULT: 'showingResult',
        TIMEOUTLOGOUT: 'idle',
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

        // disable body until result is shown
        document.body.classList.add("busy");
        
      },
    },
    // show result
    showingResult: {
      transitions: {
        SEEN: 'idle',
        TIMEOUTLOGOUT: 'idle',
        SUBMITERROR: 'selected'
      },
      onEntry: function(prevState, thisState, payload) {
        // enable body
        document.body.classList.remove("busy");

        // display body
        if (payload && payload.result) {
          const resultEl = document.getElementById('result');
          resultEl.textContent = payload.result;
        }

        // show button depending on error
        // TODO: different state for each SUBMITOK, SUBMITERROR transition
        if (payload && payload.error == true) {
            document.getElementById("resubmitBtn").hidden = false;
        } else {
            document.getElementById("seenBtn").hidden = false;
        }
      },
      onExit: function(thisState, nextState, payload) {
        // local submit error (wrong file, no input)
        // init would reset everything
        // 
        if (nextState === "selected") {

          // hide seen button
          document.getElementById("seenBtn").hidden = true;
          document.getElementById("resubmitBtn").hidden = true;

          // clean up result
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
    },
  }
}

// create Machine
const UIHandler = createMachine(UIMachine);

// register events

// on load check authentication state, as we can be redirected from callback
document.addEventListener("DOMContentLoaded", function() {
    UIHandler.transition('START');
});

// logout button
// might not exist if auth is disabled
const logoutBtn = document.getElementById("logoutBtn")
if(logoutBtn) {
  logoutBtn.addEventListener("click", () => {
    UIHandler.transition('TIMEOUTLOGOUT', {logout: true});
  });
}

// seen button
// acknowledges the result was seen
document.getElementById("seenBtn").addEventListener("click", () => {
  UIHandler.transition('SEEN');
});

// resubmit button
document.getElementById("resubmitBtn").addEventListener("click", () => {
  UIHandler.transition('SUBMITERROR');
});

// submit button
// has to be async
document.getElementById("submitBtn").addEventListener("click", async () => {
  // submit
  UIHandler.transition('SUBMIT');

  // wait for reply or handle error
  try {
    const result = await handleSubmitClick("/api/solvers/{day}/{part}");
    UIHandler.transition('SHOWRESULT', {error: false, result: result});
  } catch(error) {
    UIHandler.transition('SHOWRESULT', {error: true, result: error.message });
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
  let apiEndpoint
  try {
      apiEndpoint = fillTemplateFromSession(endpointTemplate)
  } catch (error) {
    throw Error("local error: " + error);
  }
  
  // check if something was filled
  // display error if not
  if (!file && text.trim() === "" ) {
    throw Error("local error: Please select a file or input text first");
  }
  
  // base64 encode the input, send request to API and display return value
  try {
    let base64;
    if (file) {
      base64 = await toBase64(file);
    } else {
      base64 = btoa(unescape(encodeURIComponent(text)));
    }
    
    const response = await sendToApi("POST", apiEndpoint, { input: base64 })
    return JSON.stringify(response, null, 2)
  } catch(error) {
    throw Error("backend error: " + error.message);
  }
}