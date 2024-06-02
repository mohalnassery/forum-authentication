document.querySelector(".login-form").addEventListener("submit", async (e) => {
  e.preventDefault(); // prevents the page from refreshing when form is submitted

  const username = document.querySelector("#username").value;
  const password = document.querySelector("#password").value;

  const jsonFormat = JSON.stringify({ username, password });

  try {
    const res = await fetch("/auth/login", {
      method: "POST",
      body: jsonFormat,
      headers: { "Content-Type": "application/json" },
      credentials: "include", // Ensure cookies are included
    });

    if (res.status === 409) {
      // Session already exists, invalidate the current session and retry login
      await logout();
      const retryRes = await fetch("/auth/login", {
        method: "POST",
        body: jsonFormat,
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      });

      if (!retryRes.ok) {
        const errorMessage = await retryRes.text();
        throw new Error(errorMessage);
      }

      window.location.href = "/";
    } else if (!res.ok) {
      const errorMessage = await res.text();
      throw new Error(errorMessage);
    }

    window.location.href = "/";
  } catch (error) {
    document.getElementById("error").innerHTML = error.message;
  }
});