function loadRepositories() {
  chrome.runtime.sendMessage({ type: "GET_REPOS" }, (res) => {

    if (!res || res.error) {
      document.getElementById("status").innerText =
        "Backend error: " + (res?.error || "unknown");
      return;
    }

    if (!res.repos || !Array.isArray(res.repos)) {
      document.getElementById("status").innerText =
        "No repos returned from backend";
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
