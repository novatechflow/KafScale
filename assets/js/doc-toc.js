(() => {
  const toc = document.querySelector(".doc-toc");
  const tocNav = document.querySelector(".doc-toc-nav");
  const content = document.querySelector(".doc-content");

  if (!toc || !tocNav || !content) {
    return;
  }

  const headings = Array.from(content.querySelectorAll("h2")).filter(
    (heading) => heading.textContent.trim().length > 0
  );

  if (headings.length === 0) {
    toc.style.display = "none";
    return;
  }

  const slugify = (text) =>
    text
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, "")
      .trim()
      .replace(/\s+/g, "-");

  const ensureId = (heading) => {
    if (heading.id) {
      return heading.id;
    }
    const base = slugify(heading.textContent);
    let id = base;
    let counter = 2;
    while (document.getElementById(id)) {
      id = `${base}-${counter}`;
      counter += 1;
    }
    heading.id = id;
    return id;
  };

  headings.forEach((heading) => {
    const id = ensureId(heading);
    const link = document.createElement("a");
    link.href = `#${id}`;
    link.textContent = heading.textContent;
    link.classList.add("doc-toc-link");
    tocNav.appendChild(link);
  });
})();
