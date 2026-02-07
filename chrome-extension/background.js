// let authToken = "";

// chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {

//   if (msg.type === "SET_TOKEN") {
//     authToken = msg.token;
//     sendResponse({ success: true });
//   }

//   if (msg.type === "GET_REPOS") {
//     fetch("https://git-analyser-nu1i.onrender.com/repos", {
//       headers: { Authorization: authToken }
//     })
//       .then(async res => {
//         const text = await res.text();
//         console.log("Backend response:", text);

//         try {
//           const data = JSON.parse(text);
//           sendResponse({ repos: data });
//         } catch {
//           sendResponse({ error: text });
//         }
//       })
//       .catch(err => sendResponse({ error: err.toString() }));

//     return true;
//   }
// });
// Save token persistently
chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {

  // SAVE TOKEN
  if (msg.type === "SET_TOKEN") {
    chrome.storage.local.set({ github_token: msg.token }, () => {
      sendResponse({ success: true });
    });
    return true;
  }

  // GET REPOS
  if (msg.type === "GET_REPOS") {
    chrome.storage.local.get(["github_token"], async (data) => {

      const token = data.github_token;
      if (!token) {
        sendResponse({ error: "Token missing" });
        return;
      }

      try {
        const res = await fetch("https://git-analyser-nu1i.onrender.com/repos", {
          headers: { Authorization: token }
        });

        const text = await res.text();
        console.log("Repos response:", text);

        const json = JSON.parse(text);
        sendResponse({ repos: json });

      } catch (e) {
        sendResponse({ error: e.toString() });
      }

    });
    return true;
  }

  // SYNC REPO
  if (msg.type === "SYNC_REPO") {
    chrome.storage.local.get(["github_token"], async (data) => {

      const token = data.github_token;
      const [owner, repo] = msg.repo.split("/");

      try {
        const res = await fetch(
          `https://git-analyser-nu1i.onrender.com/sync?owner=${owner}&repo=${repo}`,
          { headers: { Authorization: token } }
        );

        const text = await res.text();
        sendResponse({ message: text });

      } catch (e) {
        sendResponse({ error: e.toString() });
      }

    });
    return true;
  }

});
