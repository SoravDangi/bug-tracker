import { API_BASE_URL } from "@/config";

export interface BugLink {
  id: number;
  bugId: number;
  title: string;
  url: string;
}

export async function getBugLinks(bugId: number): Promise<BugLink[]> {
  const res = await fetch(`${API_BASE_URL}/api/bugs/${bugId}/links`);
  if (!res.ok) throw new Error("Failed to fetch links");
  return res.json();
}

export async function createBugLink(
  bugId: number,
  title: string,
  url: string
) {
  const res = await fetch(`${API_BASE_URL}/api/bugs/${bugId}/links`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, url }),
  });

  if (!res.ok) throw new Error("Failed to create link");
  return res.json();
}

export async function deleteBugLink(id: number) {
  const res = await fetch(`${API_BASE_URL}/api/bugs/links/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) throw new Error("Failed to delete link");
}
