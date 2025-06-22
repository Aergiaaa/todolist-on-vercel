function submitTodoForm(action) {
  // Get form data
  const form = document.getElementById("todo-form");
  const title = document.getElementById("title").value;
  const description = document.getElementById("description").value;

  if (!title) {
    alert("Title is required!");
    return;
  }

  let url = "/todos/";
  let formData = new FormData();
  formData.append("title", title);
  formData.append("description", description);

  if (action === "update") {
    url = "/todos/update";
    const id = document.querySelector('input[name="id"]').value;
    formData.append("id", id);
  }

  // Use HTMX to submit the form
  htmx
    .ajax("POST", url, {
      target: "#todos-container",
      swap: "innerHTML",
      values: {
        title: title,
        description: description,
        id:
          action === "update"
            ? document.querySelector('input[name="id"]').value
            : "",
      },
      headers: {
        "HX-Request": "true",
      },
    })
    .then(() => {
      // Clear form and show add button
      document.getElementById("form-container").innerHTML = "";
      document.getElementById("actions-container").style.display = "block";
    });
}
