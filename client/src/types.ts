export interface Message {
  type: "success" | "warning" | "error";
  message: string;
  description?: string;
}
