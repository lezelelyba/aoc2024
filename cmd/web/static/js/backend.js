// gets all avaialble backends
function getBackends() {
  let backends

  if (window.AppConfig && window.AppConfig.backends) {
    backends = window.AppConfig.backends;
  } else {
    backends = [];
  }

  return backends
}

// checks if backend is available
async function checkBackend(b) {

  if (!b) {
    return { backend: b, error: "backend undefined"};
  }
  
  if (typeof sendToApiDirect != "function") {
    return { backend: b, error: "call api function undefined"};
  }

  try {
    const info = await sendToApiDirect("GET", b.baseApiUrl + "/api/public/info", null);
    return { backend: b, info: info };
  } catch(error) {
    return { backend: b, error: error};
  }
}

// sets selected backend
function setBackend(backend) {
  // find the data element  
  el = document.getElementById("selected-backend");
  if (el) {
    el.dataset.name = backend.name;
  }
}

// gets selected backend
function getBackend() {

  // find the data element
  el = document.getElementById("selected-backend");
  if (!el) {
    return null;
  }

  // get backend name
  backendName = el.dataset.name;

  // find the backend
  if (window.AppConfig && window.AppConfig.backends) {
    return window.AppConfig.backends.find(b => b.name === backendName) || null;
  }

  // null if we do not have backends
  return null;
}


function updateBackendUI(backend, available, used) {

  if(!backend) {
    return;
  }

  el = document.getElementById(`backend-${backend.name}`);

  if(el) {
    el.classList.toggle("backend-used", used);
    el.classList.toggle("backend-available", (available || used));
    el.classList.toggle("backend-not-available", !available && !used);
  }
}

function setAuth(auth) {
  const configEl = document.getElementById("auth-enabled");

  if (!configEl) {
    return
  }

  if (auth == "oauth") {
    configEl.dataset.enabled = "true";
    return
  }

  configEl.dataset.enabled = "false";

  return
}
