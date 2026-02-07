let authToken = "";

chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {

  if (msg.type === "SET_TOKEN") {
    authToken = msg.token;
    sendResponse({ success: true });
  }

  if (msg.type === "GET_REPOS") {
    fetch("https://git-analyser-nu1i.onrender.com/repos", {
      headers: { Authorization: authToken }
    })
      .then(async res => {
        const text = await res.text();
        console.log("Backend response:", text);

        try {
          const data = JSON.parse(text);
          sendResponse({ repos: data });
        } catch {
          sendResponse({ error: text });
        }
      })
      .catch(err => sendResponse({ error: err.toString() }));

    return true;
  }
});
