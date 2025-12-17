import { useEffect, useState } from "react";
import { getBugLinks, createBugLink } from "@/api/buglink";

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

  return (
    <div className="mt-8">
      <h2 className="text-lg font-semibold mb-3">Reference Links</h2>

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

      <ul className="mt-4 space-y-2">
        {links.map((l) => (
          <li key={l.id}>
            🔗{" "}
            <a
              href={l.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 underline"
            >
              {l.title}
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
}
