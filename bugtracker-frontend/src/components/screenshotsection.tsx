import { useEffect, useState } from "react";
import {
  Screenshot,
  getScreenshots,
  uploadScreenshot,
  deleteScreenshot,
} from "@/api/screenshot";

import { API_BASE_URL } from "@/config";


interface Props {
  bugId: number;
}

export default function ScreenshotSection({ bugId }: Props) {
  const [screenshots, setScreenshots] = useState<Screenshot[]>([]);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  const loadScreenshots = async () => {
    try {
      const data = await getScreenshots(bugId);
      setScreenshots(data);
    } catch (err) {
      console.error(err);
      setError("Failed to load screenshots");
    }
  };

  useEffect(() => {
    loadScreenshots();
  }, [bugId]);

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files || e.target.files.length === 0) return;

    setUploading(true);
    try {
      await uploadScreenshot(bugId, e.target.files[0]);
      await loadScreenshots();
    } catch (err) {
      console.error(err);
      setError("Upload failed");
    } finally {
      setUploading(false);
      e.target.value = "";
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Delete this screenshot?")) return;

    try {
      await deleteScreenshot(id);
      setScreenshots((prev) => prev.filter((s) => s.id !== id));
    } catch (err) {
      console.error(err);
      alert("Failed to delete screenshot");
    }
  };

  return (
    <div className="mt-8">
      <h2 className="text-lg font-semibold mb-4">Screenshots</h2>

      <input
        type="file"
        accept="image/*"
        disabled={uploading}
        onChange={handleUpload}
        className="mb-4"
      />

      {error && <p className="text-red-500 mb-2">{error}</p>}

      {screenshots.length === 0 && (
        <p className="text-gray-500">No screenshots uploaded.</p>
      )}

      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        {screenshots.map((s) => (
          <div
            key={s.id}
            className="border rounded p-2 bg-white shadow"
          >
            <img
              src={`${API_BASE_URL}/${s.file_path}`}
              alt="Bug screenshot"
              className="w-full h-40 object-cover rounded cursor-pointer hover:opacity-80"
              onClick={() =>
                setPreviewUrl(`${API_BASE_URL}/${s.file_path}`)
              }
            />

            <button
              onClick={() => handleDelete(s.id)}
              className="mt-2 w-full bg-red-500 hover:bg-red-600 text-white text-sm py-1 rounded"
            >
              Delete
            </button>
          </div>
        ))}
      </div>

      {/* FULLSCREEN IMAGE PREVIEW */}
      {previewUrl && (
        <div
          className="fixed inset-0 bg-black bg-opacity-80 flex items-center justify-center z-50"
          onClick={() => setPreviewUrl(null)}
        >
          <img
            src={previewUrl}
            alt="Screenshot preview"
            className="max-w-[90%] max-h-[90%] rounded shadow-lg"
          />
        </div>
      )}
    </div>
  );
}
