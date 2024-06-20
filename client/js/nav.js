import { fetchNotifications, clearNotification, markAllNotificationsAsRead } from './notifications.js';
import { fetchUserStats } from './stats.js'; // Add this import

function createNavMenu() {
  // Inject CSS styles
  const style = document.createElement("style");
  style.innerHTML = `
      .header {
          position: fixed;
          top: 0;
          left: 0;
          right: 0;
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 10px 30px;
          background-color: #ffffff;
          box-shadow: 4px 0 5px 0 rgb(0, 0, 0, 0.3);
          z-index: 999;
      }

    .notification-icon {
      position: relative;
      cursor: pointer;
    }

    .notification-dropdown {
      display: none;
      position: absolute;
      right: 0;
      top: 100%;
      background-color: white;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
      width: 300px;
      max-height: 400px;
      overflow-y: auto;
      z-index: 1000;
    }

    .notification-dropdown.show {
      display: block;
    }

    .notification-item {
      padding: 10px;
      border-bottom: 1px solid #ddd;
      cursor: pointer;
    }

    .notification-item:last-child {
      border-bottom: none;
    }

    .notification-item:hover {
      background-color: #f5f5f5;
    }

    .mark-all-read {
      padding: 10px;
      text-align: center;
      cursor: pointer;
      background-color: #f5f5f5;
    }

    .mark-all-read:hover {
      background-color: #e0e0e0;
    }
  `;
  document.head.appendChild(style);

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

  // Add bell icon for notifications (initially hidden)
  const notificationIcon = document.createElement("div");
  notificationIcon.className = "notification-icon";
  notificationIcon.style.display = "none"; // Hide by default
  const bellIcon = document.createElement("i");
  bellIcon.className = "fa-solid fa-bell";
  notificationIcon.appendChild(bellIcon);
  userAuth.appendChild(notificationIcon);

  // Create dropdown menu for notifications
  const notificationDropdown = document.createElement("div");
  notificationDropdown.className = "notification-dropdown";
  notificationIcon.appendChild(notificationDropdown);

  // Add event listener to toggle dropdown visibility and fetch notifications
  notificationIcon.addEventListener("click", async (event) => {
    event.stopPropagation();
    notificationDropdown.classList.toggle("show");
    if (notificationDropdown.classList.contains("show")) {
      await fetchNotifications();
    }
  });

  // Hide the dropdown when clicking outside
  document.addEventListener("click", (event) => {
    if (!notificationIcon.contains(event.target)) {
      notificationDropdown.classList.remove("show");
    }
  });

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

  // Add event listener to username element to redirect to activity page
  usernameElement.addEventListener("click", () => {
    window.location.href = "/activity.html";
  });
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
        document.querySelector(".notification-icon").style.display = "inline-block"; // Show bell icon
        document.getElementById("username").textContent = data.username;
        localStorage.setItem("isLoggedIn", "true");
        localStorage.setItem("username", data.username);
        localStorage.setItem("sessionID", data.sessionID);
        fetchUserStats(); // Fetch user stats if logged in
      } else {
        isLoggedIn = false;
        document.getElementById("login-btn").style.display = "inline-block";
        document.getElementById("signup-btn").style.display = "inline-block";
        document.getElementById("logout-btn").style.display = "none";
        document.querySelector(".notification-icon").style.display = "none"; // Hide bell icon
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
      document.querySelector(".notification-icon").style.display = "none"; // Hide bell icon
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
      // Clear localStorage
      localStorage.removeItem("isLoggedIn");
      localStorage.removeItem("username");

      // Dispatch loginStatusUpdate event
      const event = new CustomEvent("loginStatusUpdate", {
        detail: { isLoggedIn: false },
      });
      window.dispatchEvent(event);

      // Redirect to home page after a short delay to ensure logout process completes
      setTimeout(() => {
        window.location.href = "/";
      }, 500); // Adjust the delay as needed
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
    document.querySelector(".notification-icon").style.display = "none"; // Hide bell icon
    document.getElementById("username").textContent = "";

    // Remove the stored session ID
    localStorage.removeItem("sessionID");
  }
});
