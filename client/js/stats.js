// Function to fetch and display user stats
export async function fetchUserStats() {
    const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";
    if (isLoggedIn) {
      try {
        const response = await fetch("/user-stats");
        if (response.ok) {
          const stats = await response.json();
          document.querySelectorAll("#user-posts").forEach((element) => {
            element.textContent = stats.posts;
          });
          document.querySelectorAll("#user-comments").forEach((element) => {
            element.textContent = stats.comments;
          });
          document.querySelectorAll("#user-likes").forEach((element) => {
            element.textContent = stats.likes;
          });
          document.querySelectorAll("#user-dislikes").forEach((element) => {
            element.textContent = stats.dislikes;
          });
        } else {
          console.error("Error fetching user stats:", response.status);
        }
      } catch (error) {
        console.error("Error fetching user stats:", error);
      }
    } else {
      // Clear user stats when the user is not logged in
      document.querySelectorAll("#user-posts").forEach((element) => {
        element.textContent = "";
      });
      document.querySelectorAll("#user-comments").forEach((element) => {
        element.textContent = "";
      });
      document.querySelectorAll("#user-likes").forEach((element) => {
        element.textContent = "";
      });
      document.querySelectorAll("#user-dislikes").forEach((element) => {
        element.textContent = "";
      });
    }
  }
  
  // Function to fetch and display all user stats
  export async function fetchAllUserStats() {
    try {
      const response = await fetch("/all-stats");
      if (!response.ok) {
        const errorMessage = await response.text();
        throw new Error(errorMessage);
      }
      const stats = await response.json();
  
      document.querySelectorAll("#total-posts").forEach((element) => {
        element.textContent = stats.totalPosts;
      });
      document.querySelectorAll("#total-comments").forEach((element) => {
        element.textContent = stats.totalComments;
      });
      document.querySelectorAll("#total-likes").forEach((element) => {
        element.textContent = stats.totalLikes;
      });
      document.querySelectorAll("#total-dislikes").forEach((element) => {
        element.textContent = stats.totalDislikes;
      });
    } catch (error) {
      console.log(error.message);
    }
  }
  
  // Function to fetch and display the leaderboard
  export async function fetchLeaderboard() {
    try {
      const response = await fetch("/leaderboard");
      if (!response.ok) {
        const errorMessage = await response.text();
        throw new Error(errorMessage);
      }
      const leaderboard = await response.json();
  
      const leaderboardContainers =
        document.querySelectorAll("#user-leaderboard");
  
      leaderboardContainers.forEach((leaderboardContainer) => {
        leaderboardContainer.innerHTML = ""; // Clear existing content
  
        leaderboard.forEach((user) => {
          const userProfile = document.createElement("div");
          userProfile.className = "user-profile";
  
          const userData = document.createElement("div");
          userData.className = "user-data";
  
          const avatar = document.createElement("div");
          avatar.className = "avatar";
          avatar.textContent = user.username.charAt(0).toUpperCase(); // Display first letter of username as avatar
          userData.appendChild(avatar);
  
          const username = document.createElement("p");
          username.textContent = user.username;
          userData.appendChild(username);
  
          userProfile.appendChild(userData);
  
          const postCount = document.createElement("p");
          postCount.textContent = user.postCount;
          userProfile.appendChild(postCount);
  
          leaderboardContainer.appendChild(userProfile);
        });
      });
    } catch (error) {
      console.log(error.message);
    }
  }
  