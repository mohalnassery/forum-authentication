function createNavMenu() {
  const header = document.createElement("div");
  header.className = "header";
  const headerLink = document.createElement("a");
  headerLink.href = "/";
  header.appendChild(headerLink);
  const logo = document.createElement("img");
  logo.src = "/assets/logo.png";
  logo.alt = "Logo";
  logo.className = "logo";
  logo.style.height = "50px";
  logo.style.width = "auto";
  headerLink.appendChild(logo);
  const userAuth = document.createElement("div");
  userAuth.className = "user-auth";
  const loginLink = document.createElement("a");
  loginLink.href = "/login";
  loginLink.className = "auth-links";
  const loginButton = document.createElement("button");
  loginButton.className = "link-buttons";
  loginButton.id = "login-btn";
  loginButton.textContent = "Login";
  loginLink.appendChild(loginButton);
  userAuth.appendChild(loginLink);
  const signupLink = document.createElement("a");
  signupLink.href = "/register";
  signupLink.className = "auth-links";
  const signupButton = document.createElement("button");
  signupButton.className = "link-buttons";
  signupButton.id = "signup-btn";
  signupButton.textContent = "SignUp";
  signupLink.appendChild(signupButton);
  userAuth.appendChild(signupLink);
  const logoutButton = document.createElement("button");
  logoutButton.id = "logout-btn";
  logoutButton.className = "link-buttons";
  logoutButton.style.display = "none";
  logoutButton.textContent = "Logout";
  logoutButton.addEventListener("click", logout);
  userAuth.appendChild(logoutButton);
  const userInfo = document.createElement("div");
  userInfo.className = "user-info";
  const usernameElement = document.createElement("p");
  usernameElement.id = "username";
  userInfo.appendChild(usernameElement);
  userAuth.appendChild(userInfo);
  header.appendChild(userAuth);
  document.body.insertBefore(header, document.body.firstChild);
}
let isLoggedIn = false; // Variable to track login status
async function updateNavMenu() {
  try {
    const response = await fetch("/auth/is-logged-in");
    if (response.ok) {
      const data = await response.json();
      if (data.status === "logged_in") {
        isLoggedIn = true;
        document.getElementById("login-btn").style.display = "none";
        document.getElementById("signup-btn").style.display = "none";
        document.getElementById("logout-btn").style.display = "inline-block";
        document.getElementById("username").textContent = data.username;
        localStorage.setItem("isLoggedIn", "true");
        localStorage.setItem("username", data.username);
        localStorage.setItem("sessionID", data.sessionID);
      } else {
        isLoggedIn = false;
        document.getElementById("login-btn").style.display = "inline-block";
        document.getElementById("signup-btn").style.display = "inline-block";
        document.getElementById("logout-btn").style.display = "none";
        document.getElementById("username").textContent = "";
        localStorage.removeItem("isLoggedIn");
        localStorage.removeItem("username");
        localStorage.removeItem("sessionID");
      }
      // Dispatch a custom event to notify the index page
      const event = new CustomEvent("loginStatusUpdate", {
        detail: { isLoggedIn },
      });
      window.dispatchEvent(event);
    } else if (response.status === 401) {
      isLoggedIn = false;
      document.getElementById("login-btn").style.display = "inline-block";
      document.getElementById("signup-btn").style.display = "inline-block";
      document.getElementById("logout-btn").style.display = "none";
      document.getElementById("username").textContent = "";
      localStorage.removeItem("isLoggedIn");
      localStorage.removeItem("username");
      localStorage.removeItem("sessionID");
      // Dispatch a custom event to notify the index page
      const event = new CustomEvent("loginStatusUpdate", {
        detail: { isLoggedIn },
      });
      window.dispatchEvent(event);
    } else {
      console.error("Error updating nav menu:", response.status);
    }
  } catch (error) {
    console.error("Error updating nav menu:", error);
  }
}
async function logout() {
  try {
    const response = await fetch("/auth/logout", { method: "POST" });
    if (response.ok) {
      window.location.href = "/";
    } else {
      console.error("Logout failed");
    }
  } catch (error) {
    console.error("Error during logout:", error);
  }
}
document.addEventListener("DOMContentLoaded", () => {
  createNavMenu();
  updateNavMenu();
});
window.addEventListener("storage", (event) => {
  if (event.key === "logout") {
    // Perform logout actions
    isLoggedIn = false;
    document.getElementById("login-btn").style.display = "inline-block";
    document.getElementById("signup-btn").style.display = "inline-block";
    document.getElementById("logout-btn").style.display = "none";
    document.getElementById("username").textContent = "";

    // Remove the stored session ID
    localStorage.removeItem("sessionID");
  }
});
