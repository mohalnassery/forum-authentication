document.addEventListener("DOMContentLoaded", () => {
  checkLoginStatus();
  createCategoryCheckboxes();
});

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
