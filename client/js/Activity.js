document.addEventListener("DOMContentLoaded", () => {
  fetchUserActivity();
});

async function fetchUserActivity() {
  try {
    const response = await fetch("/user-activity");
    if (!response.ok) {
      throw new Error("Failed to fetch user activity");
    }
    const activity = await response.json();
    displayUserActivity(activity);
  } catch (error) {
    console.error(error);
  }
}

function displayUserActivity(activity) {
  const activitySection = document.getElementById("activity-section");

  // Clear any existing content
  activitySection.innerHTML = '';

  // Ensure activity properties are initialized
  const createdPosts = activity.createdPosts || [];
  const likedPosts = activity.likedPosts || [];
  const dislikedPosts = activity.dislikedPosts || [];
  const comments = activity.comments || [];

  // Display created posts
  if (createdPosts.length > 0) {
    const createdPostsSection = document.createElement("div");
    createdPostsSection.className = "activity-subsection";
    createdPostsSection.innerHTML = "<h2>Created Posts</h2>";
    createdPosts.forEach(post => {
      const postElement = document.createElement("div");
      postElement.className = "activity-item";
      postElement.innerHTML = `<h3>${post.title}</h3><p>${post.body}</p>`;
      createdPostsSection.appendChild(postElement);
    });
    activitySection.appendChild(createdPostsSection);
  }

  // Display liked posts
  if (likedPosts.length > 0) {
    const likedPostsSection = document.createElement("div");
    likedPostsSection.className = "activity-subsection";
    likedPostsSection.innerHTML = "<h2>Liked Posts</h2>";
    likedPosts.forEach(post => {
      const postElement = document.createElement("div");
      postElement.className = "activity-item";
      postElement.innerHTML = `<h3>${post.title}</h3><p>${post.body}</p>`;
      likedPostsSection.appendChild(postElement);
    });
    activitySection.appendChild(likedPostsSection);
  }

  // Display disliked posts
  if (dislikedPosts.length > 0) {
    const dislikedPostsSection = document.createElement("div");
    dislikedPostsSection.className = "activity-subsection";
    dislikedPostsSection.innerHTML = "<h2>Disliked Posts</h2>";
    dislikedPosts.forEach(post => {
      const postElement = document.createElement("div");
      postElement.className = "activity-item";
      postElement.innerHTML = `<h3>${post.title}</h3><p>${post.body}</p>`;
      dislikedPostsSection.appendChild(postElement);
    });
    activitySection.appendChild(dislikedPostsSection);
  }

  // Display comments
  if (comments.length > 0) {
    const commentsSection = document.createElement("div");
    commentsSection.className = "activity-subsection";
    commentsSection.innerHTML = "<h2>Comments</h2>";
    comments.forEach(comment => {
      const commentElement = document.createElement("div");
      commentElement.className = "activity-item";
      commentElement.innerHTML = `<h3>Comment on: ${comment.post.title}</h3><p>${comment.body}</p>`;
      commentsSection.appendChild(commentElement);
    });
    activitySection.appendChild(commentsSection);
  }
}
