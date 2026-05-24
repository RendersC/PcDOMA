export type ClassValue =
  | string
  | number
  | null
  | false
  | undefined
  | Record<string, boolean | null | undefined>
  | ClassValue[];

/**
 * Lightweight className combiner (clsx-style) without extra dependencies.
 * Accepts strings, conditional objects, and nested arrays.
 */
export function cn(...inputs: ClassValue[]): string {
  const out: string[] = [];

  const push = (value: ClassValue) => {
    if (!value) return;
    if (typeof value === 'string' || typeof value === 'number') {
      out.push(String(value));
    } else if (Array.isArray(value)) {
      value.forEach(push);
    } else if (typeof value === 'object') {
      for (const key in value) {
        if (value[key]) out.push(key);
      }
    }
  };

  inputs.forEach(push);
  return out.join(' ');
}
