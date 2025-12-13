import { test, expect } from "@playwright/test";

/* ---------------------------
   Global setup for ALL tests
---------------------------- */
test.beforeEach(async ({ page }) => {
  page.on("console", (msg) => {
    console.log(`Browser console [${msg.type()}]: ${msg.text()}`);
  });

  await page.goto("/", { waitUntil: "domcontentloaded" });
  await page.waitForLoadState("networkidle");

  const addBugButton = page.getByRole("button", { name: "Add New Bug" });
  await expect(addBugButton).toBeVisible();
  await expect(addBugButton).toBeEnabled();
});

/* ---------------------------
   Helper: create a bug safely
---------------------------- */
async function createBug(page, title: string) {
  const addBugButton = page.getByRole("button", { name: "Add New Bug" });
  await addBugButton.click({ force: true });

  const form = page.locator('[data-testid="add-bug-form"]');
  await expect(form).toBeVisible({ timeout: 90000 });

  await page.fill('input[name="title"]', title);
  await page.fill('textarea[name="description"]', "Created by Playwright");
  await page.selectOption('select[name="priority"]', "Medium");

  await page.getByRole("button", { name: "Add Bug" }).click();

  // allow backend persistence + UI refresh
  await page.waitForTimeout(1500);
}

/* =====================================================
   TEST 1: Create and delete a bug
===================================================== */
test("Bug creation and deletion flow", async ({ page }) => {
  const uniqueTitle = `Test Bug ${Date.now()}`;

  await createBug(page, uniqueTitle);

  const bugLink = page.locator(`a:text("${uniqueTitle}")`).first();
  await expect(bugLink).toBeVisible();
  await bugLink.click();

  await expect(page.locator(`h1:text("${uniqueTitle}")`)).toBeVisible();
  await expect(page.locator("text=Created by Playwright")).toBeVisible();
  await expect(page.locator("text=Medium")).toBeVisible();

  await page.getByRole("button", { name: "Delete Bug" }).click();
  await page.getByRole("button", { name: "Confirm" }).click();

  await expect(page.locator(`a:text("${uniqueTitle}")`)).toHaveCount(0);
});

/* =====================================================
   TEST 2: Add a comment to a bug
===================================================== */
test("Adding a comment to a bug", async ({ page }) => {
  const uniqueTitle = `Bug for Comment ${Date.now()}`;
  await createBug(page, uniqueTitle);

  await page.locator(`a:text("${uniqueTitle}")`).first().click();

  const commentForm = page.locator('[data-testid="comment-form"]');
  await expect(commentForm).toBeVisible({ timeout: 60000 });

  const suffix = Date.now();
  await page.fill('[data-testid="comment-content"]', `Test comment ${suffix}`);
  await page.fill('[data-testid="comment-author"]', `User ${suffix}`);

  await page.getByRole("button", { name: "Add Comment" }).click();

  await expect(
    page.locator(`text=Test comment ${suffix}`)
  ).toBeVisible({ timeout: 30000 });
});

/* =====================================================
   TEST 3: Edit an existing bug
===================================================== */
test("Editing a bug", async ({ page }) => {
  const uniqueTitle = `Bug to Edit ${Date.now()}`;
  await createBug(page, uniqueTitle);

  await page.locator(`a:text("${uniqueTitle}")`).first().click();

  await page.getByRole("button", { name: "Edit Bug" }).click();

  const editForm = page.locator('[data-testid="edit-bug-form"]');
  await expect(editForm).toBeVisible({ timeout: 60000 });

  const editedTitle = `Edited Bug ${Date.now()}`;
  await page.fill('input[name="title"]', editedTitle);
  await page.fill(
    'textarea[name="description"]',
    "This bug has been edited"
  );
  await page.selectOption('select[name="priority"]', "High");

  await page.getByRole("button", { name: "Save Changes" }).click();

  await expect(page.locator(`h1:text("${editedTitle}")`)).toBeVisible();
  await expect(page.locator("text=This bug has been edited")).toBeVisible();
  await expect(page.locator("text=High")).toBeVisible();
});
