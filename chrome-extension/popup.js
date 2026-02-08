// // ===============================
// // LOAD REPOS FROM BACKGROUND
// // ===============================
// function loadRepositories() {
//   chrome.runtime.sendMessage({ type: "GET_REPOS" }, (res) => {

//     if (!res || res.error) {
//       document.getElementById("status").innerText =
//         "❌ " + (res?.error || "Failed to load repos");
//       return;
//     }

//     if (!res.repos || !Array.isArray(res.repos)) {
//       document.getElementById("status").innerText =
//         "❌ No repositories returned";
//       return;
//     }

//     const select = document.getElementById("repoSelect");
//     select.innerHTML = "";

//     res.repos.forEach(repo => {
//       const option = document.createElement("option");
//       option.value = `${repo.owner.login}/${repo.name}`;
//       option.textContent = `${repo.owner.login}/${repo.name}`;
//       select.appendChild(option);
//     });

//     document.getElementById("status").innerText = "✅ Repositories loaded";
//   });
// }

// // ===============================
// // SAVE TOKEN BUTTON
// // ===============================
// document.getElementById("saveToken").addEventListener("click", () => {
//   const token = document.getElementById("tokenInput").value.trim();
//   if (!token) {
//     alert("Paste GitHub token");
//     return;
//   }

//   chrome.runtime.sendMessage({ type: "SET_TOKEN", token }, () => {
//     document.getElementById("status").innerText = "✅ Token saved";
//     loadRepositories();
//   });
// });

// // ===============================
// // SYNC REPO BUTTON
// // ===============================
// document.getElementById("sync").addEventListener("click", () => {
//   const repo = document.getElementById("repoSelect").value;

//   if (!repo) {
//     alert("Select a repository first");
//     return;
//   }

//   chrome.runtime.sendMessage({ type: "SYNC_REPO", repo }, (res) => {
//     if (!res) return;

//     if (res.error) {
//       document.getElementById("status").innerText = "❌ Sync failed";
//       return;
//     }

//     document.getElementById("status").innerText = res.message;
//   });
// });
const BASE_URL = "https://git-analyser-nu1i.onrender.com";

// LOGIN BUTTON → open GitHub OAuth
document.getElementById("login").addEventListener("click", () => {
  chrome.tabs.create({
    url: `${BASE_URL}/auth/github`
  });
});

// LOGOUT
document.getElementById("logout").addEventListener("click", () => {
  chrome.storage.local.remove(["github_token"], () => {
    document.getElementById("status").innerText = "Logged out";
    document.getElementById("repoSelect").innerHTML = "";
  });
});

// Listen for token coming from backend popup window
window.addEventListener("message", (event) => {
  if (event.data.token) {
    chrome.runtime.sendMessage(
      { type: "SET_TOKEN", token: event.data.token },
      () => {
        document.getElementById("status").innerText = "✅ Logged in!";
        loadRepositories();
      }
    );
  }
});

// Load repos
function loadRepositories() {
  chrome.runtime.sendMessage({ type: "GET_REPOS" }, (res) => {
    if (!res || !res.repos) return;

    const select = document.getElementById("repoSelect");
    select.innerHTML = "";

    res.repos.forEach(repo => {
      const option = document.createElement("option");
      option.value = `${repo.owner.login}/${repo.name}`;
      option.textContent = `${repo.owner.login}/${repo.name}`;
      select.appendChild(option);
    });

    document.getElementById("status").innerText = "Repositories loaded";
  });
}

// Sync repo
document.getElementById("sync").addEventListener("click", () => {
  const repo = document.getElementById("repoSelect").value;
  if (!repo) return alert("Select repo");

  chrome.runtime.sendMessage({ type: "SYNC_REPO", repo }, (res) => {
    document.getElementById("status").innerText = res.message;
  });
});
