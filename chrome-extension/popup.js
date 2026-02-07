let authToken = "";

// Save token
document.getElementById("saveToken").addEventListener("click", () => {
  const token = document.getElementById("tokenInput").value.trim();
  if (!token) return alert("Paste token");

  chrome.runtime.sendMessage({ type: "SET_TOKEN", token });
  document.getElementById("status").innerText = "Token saved";

  loadRepositories();
});

function loadRepositories() {
  chrome.runtime.sendMessage({ type: "GET_REPOS" }, (res) => {
    if (!res || res.error) {
      document.getElementById("status").innerText = "❌ Failed to load repos";
      return;
    }

    const select = document.getElementById("repoSelect");
    select.innerHTML = "";

    res.repos.forEach(repo => {
      const option = document.createElement("option");
      option.value = `${repo.owner.login}/${repo.name}`;
      option.textContent = `${repo.owner.login}/${repo.name}`;
      select.appendChild(option);
    });

    document.getElementById("status").innerText = "✅ Repositories loaded";
  });
}

// Sync repo
document.getElementById("sync").addEventListener("click", () => {
  const repo = document.getElementById("repoSelect").value;
  if (!repo) return alert("Select repo");

  chrome.runtime.sendMessage({ type: "SYNC_REPO", repo }, (res) => {
    if (!res) return;
    document.getElementById("status").innerText = res.message;
  });
});
