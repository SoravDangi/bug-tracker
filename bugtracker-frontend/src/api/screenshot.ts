import { API_BASE_URL } from "@/config";

export interface Screenshot {
  id: number;
  bug_id: number;
  file_path: string;
  created_at: string;
}

/**
 * GET screenshots for a bug
 */
export async function getScreenshots(bugId: number): Promise<Screenshot[]> {
  const res = await fetch(
    `${API_BASE_URL}/api/bugs/${bugId}/screenshots`
  );

  if (!res.ok) {
    throw new Error("Failed to fetch screenshots");
  }

  return res.json();
}

/**
 * UPLOAD screenshot
 */
export async function uploadScreenshot(
  bugId: number,
  file: File
): Promise<void> {
  const formData = new FormData();
  formData.append("file", file);

  const res = await fetch(
    `${API_BASE_URL}/api/bugs/${bugId}/screenshots`,
    {
      method: "POST",
      body: formData,
    }
  );

  if (!res.ok) {
    throw new Error("Screenshot upload failed");
  }
}

/**
 * DELETE screenshot
 */
export async function deleteScreenshot(id: number): Promise<void> {
  const res = await fetch(
    `${API_BASE_URL}/api/screenshots/${id}`,
    {
      method: "DELETE",
    }
  );

  if (!res.ok) {
    throw new Error("Failed to delete screenshot");
  }
}
