import { fetchUserStats, fetchAllUserStats, fetchLeaderboard } from './stats.js';

document.addEventListener("DOMContentLoaded", () => {
  const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";
  if (isLoggedIn) {
    fetchUserActivity();
    fetchUserStats();
    fetchAllUserStats();
    fetchLeaderboard();
    setupTabs();
  } else {
    // Handle the case where the user is not logged in
    displayNotLoggedInMessage();
  }
});

async function fetchUserActivity() {
  try {
    const response = await fetch("/user-activity");
    if (!response.ok) {
      throw new Error("Failed to fetch user activity");
    }
    const activity = await response.json();
    console.log(activity); // Log the response to check the structure
    displayUserActivity(activity);
  } catch (error) {
    console.error(error);
  }
}

function displayUserActivity(activity) {
  const createdPostsSection = document.getElementById("created-posts-section");
  const createdCommentsSection = document.getElementById("created-comments-section");
  const likedPostsSection = document.getElementById("liked-posts-section");
  const dislikedPostsSection = document.getElementById("disliked-posts-section");
  const likedCommentsSection = document.getElementById("liked-comments-section");
  const dislikedCommentsSection = document.getElementById("disliked-comments-section");

  // Clear any existing content
  createdPostsSection.innerHTML = '';
  createdCommentsSection.innerHTML = '';
  likedPostsSection.innerHTML = '';
  dislikedPostsSection.innerHTML = '';
  likedCommentsSection.innerHTML = '';
  dislikedCommentsSection.innerHTML = '';

  // Ensure activity properties are initialized
  const createdPosts = activity.createdPosts || [];
  const likedPosts = activity.likedPosts || [];
  const dislikedPosts = activity.dislikedPosts || [];
  const comments = activity.comments || [];
  const likedComments = activity.likedComments || [];
  const dislikedComments = activity.dislikedComments || [];

  // Display created posts
  if (createdPosts.length > 0) {
    createdPosts.forEach(post => {
      const postElement = document.createElement("div");
      postElement.className = "activity-item";
      postElement.innerHTML = `<h3><a href="/post-details/${post.postID}">${post.title}</a></h3><p>${post.body}</p>`;
      createdPostsSection.appendChild(postElement);
    });
  }

  // Display liked posts
  if (likedPosts.length > 0) {
    likedPosts.forEach(post => {
      const postElement = document.createElement("div");
      postElement.className = "activity-item";
      postElement.innerHTML = `<h3><a href="/post-details/${post.postID}">${post.title}</a></h3><p>${post.body}</p>`;
      likedPostsSection.appendChild(postElement);
    });
  }

  // Display disliked posts
  if (dislikedPosts.length > 0) {
    dislikedPosts.forEach(post => {
      const postElement = document.createElement("div");
      postElement.className = "activity-item";
      postElement.innerHTML = `<h3><a href="/post-details/${post.postID}">${post.title}</a></h3><p>${post.body}</p>`;
      dislikedPostsSection.appendChild(postElement);
    });
  }

  // Display comments
  if (comments.length > 0) {
    comments.forEach(comment => {
      const commentElement = document.createElement("div");
      commentElement.className = "activity-item";
      const postTitle = comment.post ? comment.post.title : "Unknown Post"; // Handle missing post field
      commentElement.innerHTML = `<h3>Comment on: <a href="/post-details/${comment.postID}">${postTitle}</a></h3><p>${comment.body}</p>`;
      createdCommentsSection.appendChild(commentElement);
    });
  }

  // Display liked comments
  if (likedComments.length > 0) {
    likedComments.forEach(comment => {
      const commentElement = document.createElement("div");
      commentElement.className = "activity-item";
      const postTitle = comment.post ? comment.post.title : "Unknown Post"; // Handle missing post field
      commentElement.innerHTML = `<h3>Liked Comment on: <a href="/post-details/${comment.postID}">${postTitle}</a></h3><p>${comment.body}</p>`;
      likedCommentsSection.appendChild(commentElement);
    });
  }

  // Display disliked comments
  if (dislikedComments.length > 0) {
    dislikedComments.forEach(comment => {
      const commentElement = document.createElement("div");
      commentElement.className = "activity-item";
      const postTitle = comment.post ? comment.post.title : "Unknown Post"; // Handle missing post field
      commentElement.innerHTML = `<h3>Disliked Comment on: <a href="/post-details/${comment.postID}">${postTitle}</a></h3><p>${comment.body}</p>`;
      dislikedCommentsSection.appendChild(commentElement);
    });
  }
}

function setupTabs() {
  const tablinks = document.querySelectorAll(".tablinks");
  const tabcontents = document.querySelectorAll(".tabcontent");

  tablinks.forEach(tablink => {
    tablink.addEventListener("click", () => {
      const tabName = tablink.getAttribute("data-tab");

      // Hide all tab contents
      tabcontents.forEach(tabcontent => {
        tabcontent.classList.remove("active");
      });

      // Remove active class from all tab links
      tablinks.forEach(tablink => {
        tablink.classList.remove("active");
      });

      // Show the selected tab content and add active class to the clicked tab link
      document.getElementById(tabName).classList.add("active");
      tablink.classList.add("active");
    });
  });

  // Set the first tab as active by default
  if (tablinks.length > 0) {
    tablinks[0].click();
  }
}

function displayNotLoggedInMessage() {
  const mainContainer = document.querySelector(".main-container");
  mainContainer.innerHTML = "<p>You need to be logged in to view this page.</p>";
}
