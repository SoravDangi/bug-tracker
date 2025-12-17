import { API_BASE_URL } from "@/config";

const API_PATH = "/api";

export const getBugLinks = async (bugId: number) => {
  const res = await fetch(`${API_BASE_URL}${API_PATH}/bugs/${bugId}/links`);
  return res.json();
};

export const addBugLink = async (
  bugId: number,
  data: { title: string; url: string }
) => {
  const res = await fetch(
    `${API_BASE_URL}${API_PATH}/bugs/${bugId}/links`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    }
  );
  return res.json();
};
