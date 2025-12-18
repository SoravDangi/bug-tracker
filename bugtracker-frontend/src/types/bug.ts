export type Priority = 'Low' | 'Medium' | 'High';

export interface Bug {
  id: number;
  title: string;
  description: string;
  status: 'Open' | 'In Progress' | 'Resolved';
  priority: Priority;
  created_at: string;   // ✅ auto from backend
  updated_at?: string;  

}

export interface BugActions {
  deleteBug: (id: number) => Promise<void>;
}