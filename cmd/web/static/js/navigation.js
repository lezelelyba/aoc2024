/**
 * Updates session storage with currently selected day and part to be solved
 */
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