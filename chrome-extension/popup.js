// ===============================
const BASE_URL = "https://git-analyser-nu1i.onrender.com";
let authToken = "";

// Load saved token
chrome.storage.local.get(["token"], (data) => {
  if (data.token) {
    authToken = data.token;
    document.getElementById("status").innerText = "✅ Token loaded";
    loadRepositories();
  }
});

// Save token
document.getElementById("saveToken").addEventListener("click", () => {
  const token = document.getElementById("tokenInput").value.trim();
  if (!token) return alert("Paste token");

  authToken = token;
  chrome.storage.local.set({ token });
  document.getElementById("status").innerText = "✅ Token saved";

  loadRepositories();
});

// Load repos
function loadRepositories() {
  fetch(`${BASE_URL}/repos`, {
    headers: { Authorization: authToken }
  })
  .then(res => res.json())
  .then(repos => {
    const select = document.getElementById("repoSelect");
    select.innerHTML = "";

    repos.forEach(repo => {
      const option = document.createElement("option");
      option.value = `${repo.owner.login}/${repo.name}`;
      option.textContent = `${repo.owner.login}/${repo.name}`;
      select.appendChild(option);
    });

    document.getElementById("status").innerText = "✅ Repositories loaded";
  })
  .catch(() => {
    document.getElementById("status").innerText = "❌ Failed to load repositories";
  });
}

// Sync repo
document.getElementById("sync").addEventListener("click", () => {
  const repoValue = document.getElementById("repoSelect").value;
  if (!repoValue) return alert("Select repo");

  const [owner, repo] = repoValue.split("/");

  fetch(`${BASE_URL}/sync?owner=${owner}&repo=${repo}`, {
    headers: { Authorization: authToken }
  })
  .then(res => res.text())
  .then(msg => {
    document.getElementById("status").innerText = msg;
    loadHistory(repoValue);
  });
});

// Load graph
function loadHistory(repo) {
  fetch(`${BASE_URL}/history?repo=${repo}`)
  .then(res => res.json())
  .then(data => {
    const ctx = document.getElementById("chart").getContext("2d");
    const labels = data.map(d => d.time);
    const scores = data.map(d => d.score);

    new Chart(ctx, {
      type: "line",
      data: {
        labels,
        datasets: [{
          label: "Repo Activity",
          data: scores,
          borderColor: "blue",
          fill: true
        }]
      }
    });
  });
}
