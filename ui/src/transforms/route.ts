import { forEach } from "lodash";
import { matchPath } from "react-router-dom";

export function fixMenuSelected(patterns: { [props: string]: string }, pathname: string): string{

  let match = pathname
    forEach(patterns, (pattern, key) => {
      if (matchPath(key, pathname)) {
        match = pattern;
        return false
      }
    })

  return match
}