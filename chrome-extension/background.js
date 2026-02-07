let authToken = "";
let selectedRepo = "";

chrome.runtime.onMessage.addListener((msg) => {
  if (msg.type === "SET_TOKEN") authToken = msg.token;
  if (msg.type === "SET_REPO") selectedRepo = msg.repo;
});

// Auto sync every 2 minutes
setInterval(() => {
  if (!authToken || !selectedRepo) return;

  const [owner, repo] = selectedRepo.split("/");

  fetch(`https://gitsense-ooly.onrender.com/sync?owner=${owner}&repo=${repo}`, {
    headers: { Authorization: authToken }
  })
    .then(res => res.text())
    .then(msg => {
      chrome.notifications.create({
        type: "basic",
        iconUrl: "icon.png",
        title: "GitSense",
        message: msg
      });
    });
}, 120000);
