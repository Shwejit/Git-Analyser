let authToken = "";
let selectedRepo = "";

chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {

  // Save token
  if (msg.type === "SET_TOKEN") {
    authToken = msg.token;
    sendResponse({ success: true });
  }

  // Fetch repos (NO CORS HERE 🎉)
  if (msg.type === "GET_REPOS") {
    fetch("https://git-analyser-nu1i.onrender.com/repos", {
      headers: { Authorization: authToken }
    })
      .then(res => res.json())
      .then(data => sendResponse({ repos: data }))
      .catch(() => sendResponse({ error: true }));

    return true; // keep channel open
  }

  // Sync repo
  if (msg.type === "SYNC_REPO") {
    const [owner, repo] = msg.repo.split("/");

    fetch(`https://git-analyser-nu1i.onrender.com/sync?owner=${owner}&repo=${repo}`, {
      headers: { Authorization: authToken }
    })
      .then(res => res.text())
      .then(msg => sendResponse({ message: msg }))
      .catch(() => sendResponse({ error: true }));

    return true;
  }

});
