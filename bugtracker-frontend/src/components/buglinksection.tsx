import { useEffect, useState } from "react";
import {
  getBugLinks,
  createBugLink,
  deleteBugLink,
} from "@/api/buglink";

interface BugLink {
  id: number;
  title: string;
  url: string;
}

export default function BugLinksSection({ bugId }: { bugId: number }) {
  const [links, setLinks] = useState<BugLink[]>([]);
  const [title, setTitle] = useState("");
  const [url, setUrl] = useState("");

  const load = async () => {
    const data = await getBugLinks(bugId);
    setLinks(Array.isArray(data) ? data : []);
  };

  useEffect(() => {
    load();
  }, [bugId]);

  const submit = async () => {
    if (!title || !url) return;
    await createBugLink(bugId, title, url);
    setTitle("");
    setUrl("");
    load();
  };

  const remove = async (id: number) => {
    if (!confirm("Delete this link?")) return;
    await deleteBugLink(id);
    load();
  };

  return (
    <div className="mt-8">
      <h2 className="text-lg font-semibold mb-3">Reference Links</h2>

      <div className="flex mb-4">
        <input
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="border p-2 mr-2 rounded"
        />

        <input
          placeholder="https://..."
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          className="border p-2 mr-2 rounded w-96"
        />

        <button
          onClick={submit}
          className="bg-blue-500 text-white px-3 py-2 rounded"
        >
          Add
        </button>
      </div>

      <ul className="mt-4 space-y-2">
        {links.map((l) => (
          <li key={l.id} className="flex items-center gap-3">
            <span>🔗</span>
            <a
              href={l.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 underline flex-1"
            >
              {l.title}
            </a>
            <button
              onClick={() => remove(l.id)}
              className="text-red-600 border px-2 py-1 rounded hover:bg-red-100"
            >
              Delete
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
