async function sendToApi(endpoint, payload) {
  console.log("Sending to:", endpoint);

  try {
    const res = await fetch(endpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(payload)
    });

    const data = await res.json();
    console.log("Response:", data);
    return data;
  } catch (err) {
    console.error("Error:", err);
    return { error: err.message };
  }
}

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