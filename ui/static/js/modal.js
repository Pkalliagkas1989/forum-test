document.addEventListener("DOMContentLoaded", () => {
  const createBtn = document.getElementById("create-post-btn");
  const modal = document.getElementById("post-modal");
  const closeBtn = document.querySelector(".close-btn");
  const submitBtn = document.getElementById("submit-post");
  const postTitle = document.getElementById("post-title");
  const postBody = document.getElementById("post-body");
  const categorySelect = document.getElementById("post-category");

  createBtn.onclick = () => {
    modal.classList.remove("hidden");
  };

  closeBtn.onclick = () => {
    modal.classList.add("hidden");
  };

  // Populate categories
  fetch("http://localhost:8080/forum/api/categories")
    .then((res) => res.json())
    .then((categories) => {
      categorySelect.innerHTML =
        '<option value="" disabled selected>Select category</option>';
      categories.forEach((c) => {
        const opt = document.createElement("option");
        opt.value = c.id;
        opt.textContent = c.name;
        categorySelect.appendChild(opt);
      });
    })
    .catch((err) => console.error("Failed to load categories", err));

  submitBtn.onclick = async () => {
    const title = postTitle.value.trim();
    const content = postBody.value.trim();
    const categoryId = parseInt(categorySelect.value, 10);

    if (!title || !content || isNaN(categoryId)) return;

    try {
      const res = await fetch("http://localhost:8080/forum/api/posts", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          title,
          content,
          category_id: categoryId,
        }),
      });
      if (!res.ok) throw new Error("Failed to create post");
      postTitle.value = "";
      postBody.value = "";
      modal.classList.add("hidden");
      window.location.reload();
    } catch (err) {
      console.error("Create post failed", err);
    }
  };
});
