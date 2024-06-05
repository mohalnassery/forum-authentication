document.addEventListener("DOMContentLoaded", () => {
  checkLoginStatus();
  createCategoryCheckboxes();
  validateImageUpload();
});

function validateImageUpload() {
  const imageInput = document.getElementById("image");
  imageInput.addEventListener("change", () => {
    const file = imageInput.files[0];
    const allowedTypes = ["image/jpeg", "image/png", "image/gif"];
    const maxSize = 20 * 1024 * 1024; // 20 MB

    if (file) {
      if (!allowedTypes.includes(file.type)) {
        alert("Only JPEG, PNG, and GIF file types are allowed.");
        imageInput.value = "";
      } else if (file.size > maxSize) {
        alert("File size exceeds the maximum limit of 20 MB.");
        imageInput.value = "";
      }
    }
  });
}
// fetchCategories.js
async function fetchCategories() {
  try {
    const response = await fetch("/categories");
    if (!response.ok) {
      throw new Error("Failed to fetch categories");
    }
    const categories = await response.json();
    return categories;
  } catch (error) {
    console.error(error);
    return [];
  }
}

async function createCategoryCheckboxes() {
  const categories = await fetchCategories();
  const container = document.querySelector(".tag-input-container");
  container.innerHTML = "";

  categories.forEach((category, idx) => {
    const label = document.createElement("label");
    const checkbox = document.createElement("input");
    checkbox.type = "checkbox";
    checkbox.name = "options";
    checkbox.value = category;
    label.appendChild(checkbox);
    label.appendChild(document.createTextNode(" " + category));
    container.appendChild(label);
  });
}

// function showPopupMessage(message) {
//   const popup = document.createElement("div");
//   popup.classList.add("popup-message");
//   popup.textContent = message;
//   document.body.appendChild(popup);

//   setTimeout(() => {
//     popup.remove();
//   }, 3000);
// }

async function checkLoginStatus() {
  const isLoggedIn = localStorage.getItem("isLoggedIn") === "true";
  if (!isLoggedIn) {
    window.location.href = "/login";
  }
}


async function createPost(event) {
  event.preventDefault();

  const form = event.target;
  const formData = new FormData(form);

  // Perform client-side validation
  const imageInput = document.getElementById("image");
  const file = imageInput.files[0];
  const allowedTypes = ["image/jpeg", "image/png", "image/gif"];
  const maxSize = 20 * 1024 * 1024; // 20 MB

  if (file) {
    if (!allowedTypes.includes(file.type)) {
      alert("Only JPEG, PNG, and GIF file types are allowed.");
      return;
    } else if (file.size > maxSize) {
      alert("File size exceeds the maximum limit of 20 MB.");
      return;
    }
  }

  try {
    const response = await fetch("/posts", {
      method: "POST",
      body: formData,
    });

    if (response.ok) {
      // Post created successfully
      window.location.href = "/";
    } else {
      // Handle error response
      console.error("Failed to create post");
    }
  } catch (error) {
    console.error("Error creating post:", error);
  }
}