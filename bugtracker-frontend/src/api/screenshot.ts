import { API_BASE_URL } from "@/config";

export interface Screenshot {
  id: number;
  bug_id: number;
  file_path: string;
}

const API_PATH = "/api";

export const getScreenshots = async (bugId: number): Promise<Screenshot[]> => {
  const res = await fetch(
    `${API_BASE_URL}${API_PATH}/bugs/${bugId}/screenshots`
  );

  if (!res.ok) {
    throw new Error("Failed to fetch screenshots");
  }

  const data = await res.json();
  return Array.isArray(data) ? data : [];
};

export const uploadScreenshot = async (
  bugId: number,
  file: File
) => {
  const formData = new FormData();
  formData.append("file", file);

  const res = await fetch(
    `${API_BASE_URL}${API_PATH}/bugs/${bugId}/screenshots`,
    {
      method: "POST",
      body: formData, // ❗ DO NOT SET HEADERS
    }
  );

  if (!res.ok) {
    throw new Error("Upload failed");
  }

  return await res.json();
};

export const deleteScreenshot = async (id: number) => {
  const res = await fetch(
    `${API_BASE_URL}${API_PATH}/screenshots/${id}`,
    {
      method: "DELETE",
    }
  );

  if (!res.ok) {
    throw new Error("Delete failed");
  }
};
