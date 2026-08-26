import { writeFileSync } from "node:fs";
writeFileSync("dist.js", "export function friendlyError(error) { return error?.error || '请求失败，请稍后重试'; }\n");
