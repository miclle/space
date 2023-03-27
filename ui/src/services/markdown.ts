import { POST } from "./lib/http";

export function preview(content: string): Promise<string> {
  return POST(`/markdown/preview`, { content })
}
