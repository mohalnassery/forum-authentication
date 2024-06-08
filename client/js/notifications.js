async function fetchNotifications() {
  try {
    const response = await fetch("/notifications");
    if (response.ok) {
      const notifications = await response.json();
      displayNotifications(notifications);
    } else {
      console.error("Error fetching notifications:", response.status);
      displayNotifications([]); // Display "NO NEW NOTIFICATION" if there's an error
    }
  } catch (error) {
    console.error("Error fetching notifications:", error);
    displayNotifications([]); // Display "NO NEW NOTIFICATION" if there's an error
  }
}

function displayNotifications(notifications) {
  const notificationDropdown = document.querySelector(".notification-dropdown");
  notificationDropdown.innerHTML = ""; // Clear existing notifications

  // Filter out read notifications
  const unreadNotifications = notifications.filter(notification => !notification.is_read);

  if (unreadNotifications.length > 0) {
    unreadNotifications.forEach((notification) => {
      const notificationItem = document.createElement("div");
      notificationItem.className = "notification-item";
      notificationItem.textContent = notification.message;
      notificationItem.addEventListener("click", () => clearNotification(notification.id));
      notificationDropdown.appendChild(notificationItem);
    });

    // Add "Mark all as read" option
    const markAllRead = document.createElement("div");
    markAllRead.className = "mark-all-read";
    markAllRead.textContent = "Mark all as read";
    markAllRead.addEventListener("click", markAllNotificationsAsRead);
    notificationDropdown.appendChild(markAllRead);
  } else {
    const noNotifications = document.createElement("div");
    noNotifications.className = "notification-item";
    noNotifications.textContent = "NO NEW NOTIFICATION";
    notificationDropdown.appendChild(noNotifications);
  }
}

async function clearNotification(notificationId) {
  try {
    const response = await fetch(`/notifications/clear/${notificationId}`, {
      method: "POST",
    });
    if (response.ok) {
      const data = await response.json();
      const postID = data.postID;
      if (postID) {
        window.location.href = `/post-details/${postID}`; // Redirect to the post details page
      } else {
        console.error("Post ID not found in the response");
      }
    } else {
      console.error("Error clearing notification:", response.status);
    }
  } catch (error) {
    console.error("Error clearing notification:", error);
  }
}


async function markAllNotificationsAsRead() {
  try {
    const response = await fetch("/notifications/mark-all-read", {
      method: "POST",
    });
    if (response.ok) {
      fetchNotifications(); // Refresh notifications
    } else {
      console.error("Error marking all notifications as read:", response.status);
    }
  } catch (error) {
    console.error("Error marking all notifications as read:", error);
  }
}

// Export functions to be used in other files
export { fetchNotifications, clearNotification, markAllNotificationsAsRead };
