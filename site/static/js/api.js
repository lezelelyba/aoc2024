/**
 * Sends an authenticated POST request to the API
 * 
 * includes Authorization header with the stored access token
 * 
 * @param {string} apiEndpoint - Target API Endpoint URL
 * @param {Object} payload - payload to be serialized as JSON
 * @returns {Promise<object>} Promise resolving to the parsed API response
 */
async function sendToApi(method, apiEndpoint, payload) {
  const accessToken = sessionStorage.getItem("accessToken");

  if (window.AppConfig && window.AppConfig.baseApiUrl) {
    apiEndpoint = window.AppConfig.baseApiUrl + apiEndpoint
  }
  
  console.log("Sending to:", apiEndpoint);

  const headers = {"Content-Type": "application/json"};
  if (accessToken) {
    headers["Authorization"] = "Bearer " + accessToken
  }

  try {
    // sends authorized request
    let res;

    if (method == "GET") {
      res = await fetch(apiEndpoint, {
        method: method,
        headers,
      });
    } else if (method == "POST") {
      res = await fetch(apiEndpoint, {
        method: method,
        headers,
        body: JSON.stringify(payload)
      });
    }

    // attempt to parse response
    const data = await res.json();
    console.log("Response:", data);
    return data;
  } catch (err) {
    console.error("Error:", err);
    throw err
  }
}

/**
 * Encodes content of a file into base64
 *  
 * @param {string} file - filepath
 * @returns {Promise<string>} Promise that resolves with Base64-encoded file content
 */
function toBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => {
      const base64 = reader.result.split(',')[1]; // strip the "data:...," part
      resolve(base64);
    }
    reader.onerror = reject;
  });
}