const loginForm = document.getElementById('loginForm');
const submitButton = document.getElementById('submitButton');
const loginSpinner = document.getElementById('loginSpinner');
const errorContainer = document.getElementById('errorContainer');

loginForm.addEventListener("submit", (e) => {
  e.preventDefault();
  errorContainer.innerHTML = "";
  errorContainer.classList.add("invisible");
  submitButton.setAttribute("disabled", "disabled");
  loginSpinner.classList.remove("invisible");
  const username = document.getElementById("username");
  const pass = document.getElementById("password");
  
  fetch("/api/login", {
    method: "POST",
    body: JSON.stringify({
      username: username.value,
      password: pass.value
    }),
    headers: {
      "Content-type": "application/json; charset=UTF-8"
    }
  })
  .then((response) => {
      if (response.status != 200) {
        console.log("Bad login");
        const error = document.createElement("p");
        error.textContent = "Failed to login. Please try again"
        errorContainer.replaceChildren(error);
        errorContainer.classList.remove("invisible");
        loginSpinner.classList.add("invisible");
        submitButton.removeAttribute("disabled");
        return;
      }
      return response.json()
    })
  .then((json) => {
      loginSpinner.classList.add("invisible");
      submitButton.removeAttribute("disabled");
      console.log(json);
  })
});
