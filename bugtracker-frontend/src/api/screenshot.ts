import { API_BASE_URL } from "@/config";

export interface Screenshot {
  id: number;
  bug_id: number;
  file_name: string;
  file_path: string;
  created_at: string;
}

// Upload screenshot
export async function uploadScreenshot(
  bugId: number,
  file: File
): Promise<Screenshot> {
  const formData = new FormData();
  formData.append("screenshot", file);

  const response = await fetch(
    `${API_BASE_URL}/api/bugs/${bugId}/screenshots`,
    {
      method: "POST",
      body: formData,
    }
  );

  if (!response.ok) {
    throw new Error("Failed to upload screenshot");
  }

  return response.json();
}

// Fetch screenshots
export async function getScreenshots(
  bugId: number
): Promise<Screenshot[]> {
  const response = await fetch(
    `${API_BASE_URL}/api/bugs/${bugId}/screenshots`
  );

  if (!response.ok) {
    throw new Error("Failed to fetch screenshots");
  }

  return response.json();
}

// Delete screenshot
export async function deleteScreenshot(id: number): Promise<void> {
  const response = await fetch(
    `${API_BASE_URL}/api/screenshots/${id}`,
    { method: "DELETE" }
  );

  if (!response.ok) {
    throw new Error("Failed to delete screenshot");
  }
}
