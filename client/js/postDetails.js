const maxCharacters = 500;

document.addEventListener("DOMContentLoaded", () => {
  fetchPostsAndStatus();
});

// Event listener for login status update
window.addEventListener("loginStatusUpdate", (event) => {
  const { isLoggedIn } = event.detail;
  if (isLoggedIn) {
    enableInteractions();
  } else {
    disableInteractions();
  }
});

async function fetchPostsAndStatus() {
  try {
    await fetchPostDetails(getPostIdFromURL());
    const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";
    if (isLoggedIn) {
      enableInteractions();
    } else {
      disableInteractions();
    }
  } catch (error) {
    console.error("Error fetching post details:", error);
    disableInteractions();
  }
}

function enableInteractions() {
  const submitCommentBtn = document.getElementById("submit-comment");
  const likeButton = document.getElementById("like-btn");
  const dislikeButton = document.getElementById("dislike-btn");
  const commentInput = document.getElementById("comment-input");
  submitCommentBtn.removeAttribute("disabled");
  likeButton.disabled = false;
  dislikeButton.disabled = false;
  likeButton.classList.remove("disabled-button");
  dislikeButton.classList.remove("disabled-button");
  commentInput.removeAttribute("disabled");

  // Event Listeners

  submitCommentBtn.addEventListener("click", submitComment);
  likeButton.addEventListener("click", likePost);
  dislikeButton.addEventListener("click", dislikePost);
  const charCount = document.getElementById("char-count");

  commentInput.addEventListener("input", function () {
    const remainingChars =
      maxCharacters -
      commentInput.textContent
        .replaceAll("<div><br></div>", "\n")
        .replaceAll("<div>", "\n")
        .replaceAll("</div>", "")
        .replaceAll("<br>", "\n")
        .replaceAll("&nbsp;"," ")
        .length;

    charCount.textContent = remainingChars + "/" + maxCharacters;
    submitCommentBtn.disabled = remainingChars < 0;
  });
}

function disableInteractions() {
  const commentInput = document.getElementById("comment-input");
  commentInput.setAttribute("contenteditable", "false");
  commentInput.setAttribute(
    "placeholder",
    "Please login to be able to comment..."
  );
  document.getElementById("submit-comment").setAttribute("disabled", "");
  const likeButton = document.getElementById("like-btn");
  const dislikeButton = document.getElementById("dislike-btn");
  likeButton.disabled = true;
  dislikeButton.disabled = true;
  likeButton.classList.add("disabled-button");
  dislikeButton.classList.add("disabled-button");
  likeButton.classList.remove("liked");
  dislikeButton.classList.remove("disliked");
}

function getPostIdFromURL() {
  const urlParts = window.location.pathname.split("/");
  return urlParts[urlParts.length - 1];
}

async function fetchPostDetails(postId) {
  try {
    const response = await fetch(`/posts/${postId}`);
    if (!response.ok) {
      const errorMessage = await response.text();
      throw new Error(
        `Failed to fetch post details. Status: ${response.status}. Message: ${errorMessage}`
      );
    }
    const data = await response.json();
    displayPostDetails(data.post);
    displayComments(data.comments);
  } catch (error) {
    window.location.href = "/400";
    console.error("Error fetching post details:", error);
  }
}

function displayPostDetails(post) {
  document.getElementById("post-author-avatar").textContent = post.author
    .charAt(0)
    .toUpperCase();
  document.getElementById("post-author").innerHTML =
    post.author +
    ` <span>• ${new Date(post.creationdate).toLocaleDateString()}</span>`;
  document.getElementById("post-comments").textContent = post.commentCount;
  document.getElementById("post-title").textContent = post.title;
  document.getElementById("post-body").innerHTML = post.body;
  document.getElementById("post-likes").textContent = post.likes;
  document.getElementById("post-dislikes").textContent = post.dislikes;
  const likeButton = document.getElementById("like-btn");
  const dislikeButton = document.getElementById("dislike-btn");
  if (isLoggedIn) {
    likeButton.classList.toggle("liked", post.userLiked);
    dislikeButton.classList.toggle("disliked", post.userDisliked);
  } else {
    likeButton.classList.remove("liked");
    dislikeButton.classList.remove("disliked");
  }

  const tagsContainer = document.getElementById("post-tags");
  if (tagsContainer) {
    tagsContainer.innerHTML = ""; // Clear existing tags

    post.categories.forEach((category) => {
      const tag = document.createElement("div");
      tag.textContent = category;
      tagsContainer.appendChild(tag);
    });
  }

  if (post.isAuthor) {
    const deletePostButton = document.createElement("i");
    deletePostButton.className = "fa-solid fa-trash";
    deletePostButton.addEventListener("click", deletePost);
    const postActions = document.getElementById("post-actions");
    if (postActions) {
      // Check if the delete post button already exists
      const existingDeleteButton = postActions.querySelector("i.fa-trash");
      if (!existingDeleteButton) {
        // Append the delete post button only if it doesn't exist
        postActions.appendChild(deletePostButton);
      }
    } else {
      console.log("Post actions container not found");
    }
  }
}

function displayComments(comments) {
  const commentList = document.getElementById("comment-section");
  commentList.innerHTML = "";

  if (comments && comments.length > 0) {
    comments.forEach((comment) => {
      const commentCard = document.createElement("div");
      commentCard.classList.add("comment");

      const commentHeader = document.createElement("div");
      commentHeader.classList.add("info");

      const avatar = document.createElement("div");
      avatar.classList.add("avatar");
      avatar.textContent = comment.author.charAt(0).toUpperCase();
      commentHeader.appendChild(avatar);

      const author = document.createElement("p");
      author.innerHTML =
        comment.author +
        ` <span>• ${new Date(
          comment.creationDate
        ).toLocaleDateString()}</span>`;
      commentHeader.appendChild(author);

      commentCard.appendChild(commentHeader);
      const commentText = document.createElement("p");
      commentText.classList.add("comment-data");
      commentText.innerHTML = comment.body;
      commentCard.appendChild(commentText);
      const commentActions = document.createElement("div");
      commentActions.classList.add("like-dislike");

      const likes = document.createElement("div");
      const likeCount = document.createElement("p");
      likeCount.classList.add("comment-likes");
      likeCount.innerText = comment.likes;
      likes.appendChild(likeCount);

      const likeButton = document.createElement("i");
      likeButton.id = `like-comment-${comment.id}`;
      likeButton.className = "fa-solid fa-thumbs-up";
      likeButton.addEventListener("click", () => likeComment(comment.id));
      if (isLoggedIn) {
        likeButton.disabled = false;
        likeButton.classList.remove("disabled-button");
      } else {
        likeButton.disabled = true;
        likeButton.classList.add("disabled-button");
      }
      likeButton.classList.toggle("liked", comment.userLiked);
      likes.appendChild(likeButton);

      commentActions.appendChild(likes);

      const dislikes = document.createElement("div");
      const dislikeCount = document.createElement("p");
      dislikeCount.classList.add("comment-dislikes");
      dislikeCount.innerText = comment.dislikes;
      dislikes.appendChild(dislikeCount);

      const dislikeButton = document.createElement("i");
      dislikeButton.id = `dislike-comment-${comment.id}`;
      dislikeButton.className = "fa-solid fa-thumbs-down";
      dislikeButton.addEventListener("click", () => dislikeComment(comment.id));
      if (isLoggedIn) {
        dislikeButton.disabled = false;
        dislikeButton.classList.remove("disabled-button");
      } else {
        dislikeButton.disabled = true;
        dislikeButton.classList.add("disabled-button");
      }
      dislikeButton.classList.toggle("disliked", comment.userDisliked);
      dislikes.appendChild(dislikeButton);

      commentActions.appendChild(dislikes);

      commentCard.appendChild(commentActions);

      if (comment.isAuthor) {
        const deleteCommentButton = document.createElement("i");
        deleteCommentButton.className = "fa-solid fa-trash";
        deleteCommentButton.addEventListener("click", () =>
          deleteComment(comment.id)
        );
        commentCard.appendChild(deleteCommentButton);
      }

      commentList.appendChild(commentCard);
    });
  } else {
    const noCommentsMessage = document.createElement("p");
    noCommentsMessage.textContent = "No comments available.";
    commentList.appendChild(noCommentsMessage);
  }

  // Clear the comment input field
  document.getElementById("comment-input").value = "";
}

async function submitComment() {
  const commentInput = document.getElementById("comment-input");
  const errorMsg = document.getElementById("error")
  const commentText = commentInput.innerHTML
    .replaceAll("<div><br></div>", "\n")
    .replaceAll("<div>", "\n")
    .replaceAll("</div>", "")
    .replaceAll("<br>", "\n")
    .replaceAll("&nbsp;"," ")
    .trim();
  if (commentText !== "") {
    errorMsg.innerText = ""
    const postId = getPostIdFromURL();
    try {
      const response = await fetch(`/posts/${postId}/comments`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ body: commentText }),
      });

      if (response.ok) {
        commentInput.textContent = ""; // Clear the comment input
        fetchPostDetails(postId); // Refresh the post details and comments
        document.getElementById("char-count").textContent =
          maxCharacters + "/" + maxCharacters; // Refresh comment char count
      } else {
        const errorMessage = await response.text();
        console.error("Error submitting comment:", errorMessage);
      }
    } catch (error) {
      console.error("Error submitting comment:", error);
    }
  } else {
    errorMsg.innerText = "Comment can not be empty"
  }
}

async function likePost() {
  if (!isLoggedIn) {
    console.log("User is not logged in. Cannot like post.");
    return;
  }

  const postId = getPostIdFromURL();

  try {
    const response = await fetch(`/posts/${postId}/like`, {
      method: "POST",
    });

    if (response.ok) {
      fetchPostDetails(postId); // Refresh the post details
    } else {
      const errorMessage = await response.text();
      console.error("Error liking post:", errorMessage);
    }
  } catch (error) {
    console.error("Error liking post:", error);
  }
}

async function dislikePost() {
  if (!isLoggedIn) {
    console.log("User is not logged in. Cannot dislike post.");
    return;
  }

  const postId = getPostIdFromURL();

  try {
    const response = await fetch(`/posts/${postId}/dislike`, {
      method: "POST",
    });

    if (response.ok) {
      fetchPostDetails(postId); // Refresh the post details
    } else {
      const errorMessage = await response.text();
      console.error("Error disliking post:", errorMessage);
    }
  } catch (error) {
    console.error("Error disliking post:", error);
  }
}

async function likeComment(commentId) {
  if (!isLoggedIn) {
    console.log("User is not logged in. Cannot like comment.");
    return;
  }

  try {
    const response = await fetch(`/comments/${commentId}/like`, {
      method: "POST",
    });

    if (response.ok) {
      fetchPostDetails(getPostIdFromURL()); // Refresh the post details
    } else {
      const errorMessage = await response.text();
      console.error("Error liking comment:", errorMessage);
    }
  } catch (error) {
    console.error("Error liking comment:", error);
  }
}

async function dislikeComment(commentId) {
  if (!isLoggedIn) {
    console.log("User is not logged in. Cannot dislike comment.");
    return;
  }

  try {
    const response = await fetch(`/comments/${commentId}/dislike`, {
      method: "POST",
    });

    if (response.ok) {
      fetchPostDetails(getPostIdFromURL()); // Refresh the post details
    } else {
      const errorMessage = await response.text();
      console.error("Error disliking comment:", errorMessage);
    }
  } catch (error) {
    console.error("Error disliking comment:", error);
  }
}

async function deletePost() {
  const postId = getPostIdFromURL();

  try {
    const response = await fetch(`/posts/${postId}`, {
      method: "DELETE",
    });

    if (response.ok) {
      window.location.href = "/"; // Redirect to the home page after successful deletion
    } else {
      const errorMessage = await response.text();
      console.error("Error deleting post:", errorMessage);
    }
  } catch (error) {
    console.error("Error deleting post:", error);
  }
}

async function deleteComment(commentId) {
  try {
    const response = await fetch(`/comments/${commentId}`, {
      method: "DELETE",
    });

    if (response.ok) {
      fetchPostDetails(getPostIdFromURL()); // Refresh the post details after successful deletion
    } else {
      const errorMessage = await response.text();
      console.error("Error deleting comment:", errorMessage);
    }
  } catch (error) {
    console.error("Error deleting comment:", error);
  }
}
